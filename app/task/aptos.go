package task

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"github.com/smallnest/chanx"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

type aptos struct {
	versionChunkSize       int64
	versionConfirmedOffset int64
	lastVersion            int64
	versionQueue           *chanx.UnboundedChan[version]
	client                 *http.Client
	earliestBlockTs        atomic.Int64
}

type version struct {
	Start int64
	Limit int64
}

var apt aptos

type aptEvent struct {
	Type    string
	Action  string
	Amount  decimal.Decimal
	Address string
}

type aptAmount struct {
	Amount string
	Type   model.TradeType
}

func init() {
	apt = newAptos()
	Register(Task{Callback: apt.versionDispatch})
	Register(Task{Callback: apt.syncVersionForward, Duration: time.Second * 3})
	Register(Task{Callback: apt.tradeConfirmHandle, Duration: time.Second * 5})
}

func newAptos() aptos {
	return aptos{
		versionChunkSize:       100, // 目前好像最大就只能100
		versionConfirmedOffset: 1000,
		lastVersion:            0,
		versionQueue:           chanx.NewUnboundedChan[version](context.Background(), 30),
		client:                 utils.NewHttpClient(),
		earliestBlockTs:        atomic.Int64{},
	}
}

func (a *aptos) syncVersionForward(ctx context.Context) {
	if syncBreak(conf.Aptos, a.versionQueue.Len()) {

		return
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", model.Endpoint(conf.Aptos)+"/v1", nil)
	resp, err := a.client.Do(req)
	if err != nil {
		log.Task.Warn("aptos syncVersionForward Error sending request:", err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Task.Warn("aptos syncVersionForward Error response status code:", resp.StatusCode)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Task.Warn("aptos syncVersionForward Error reading response body:", err)

		return
	}

	now := gjson.GetBytes(body, "ledger_version").Int()
	if now <= 0 {
		log.Task.Warn("syncVersionForward Error: invalid ledger_version:", now)

		return
	}

	if a.lastVersion == 0 {
		a.syncVersionBackward(ctx, now)
	}

	if now-a.lastVersion > 10000 {
		a.lastVersion = now - a.versionChunkSize
	}

	var sub = now - a.lastVersion
	if sub <= a.versionChunkSize {
		a.versionQueue.In <- version{Start: a.lastVersion, Limit: sub}
	} else {
		chunks := (sub + a.versionChunkSize - 1) / a.versionChunkSize
		for i := int64(0); i < chunks; i++ {
			limit := a.versionChunkSize
			start := a.lastVersion + a.versionChunkSize*i
			if i == chunks-1 {
				limit = sub % a.versionChunkSize
				if limit == 0 {
					limit = a.versionChunkSize
				}
			}

			a.versionQueue.In <- version{Start: start, Limit: limit}
		}
	}

	a.lastVersion = now
}

func (a *aptos) syncVersionBackward(ctx context.Context, now int64) {
	var o model.Order
	trade := model.GetNetworkTrades(conf.Aptos)
	model.Db.Model(&model.Order{}).Where("status = ? and trade_type in (?)",
		model.OrderStatusWaiting, trade).Order("created_at asc").Limit(1).Find(&o)
	if o.ID == 0 {
		return
	}

	start := o.CreatedAt.Time().Unix()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if syncBreak(conf.Aptos, a.versionQueue.Len()) {
					return
				}

				if a.earliestBlockTs.Load() != 0 && a.earliestBlockTs.Load() < start {
					return
				}

				a.versionQueue.In <- version{
					Start: now - a.versionChunkSize,
					Limit: a.versionChunkSize,
				}

				now = now - a.versionChunkSize
			}
		}
	}()
}

func (a *aptos) versionDispatch(ctx context.Context) {
	p, err := ants.NewPoolWithFunc(3, a.versionParse)
	if err != nil {
		log.Task.Warn("aptos versionDispatch Error:", err)

		return
	}

	defer p.Release()

	for {
		select {
		case n := <-a.versionQueue.Out:
			if err := p.Invoke(n); err != nil {
				a.versionQueue.In <- n
				log.Task.Warn("versionDispatch Error invoking process slot:", err)
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Task.Warn("versionDispatch context done:", err)
			}

			return
		}
	}
}

// 由于 aptos 网络特性，交易数据中不会显示存在交易转账 from => to 的对应关系，
// 所以目前此解析函数存在大量循环嵌套解析，逻辑较为复杂，希望未来有更好的方式进行解析 慢慢优化
func (a *aptos) versionParse(n any) {
	p := n.(version)

	var net = conf.Aptos
	var url = fmt.Sprintf("%sv1/transactions?start=%d&limit=%d", model.Endpoint(conf.Aptos), p.Start, p.Limit)

	conf.RecordSuccess(net)
	resp, err := a.client.Get(url)
	if err != nil {
		conf.RecordFailure(net)
		log.Task.Warn("versionParse Error sending request:", err)

		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		conf.RecordFailure(net)
		log.Task.Warn("versionParse Error response status code:", resp.StatusCode)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		conf.RecordFailure(net)
		a.versionQueue.In <- p
		log.Task.Warn("versionParse Error reading response body:", err)

		return
	}

	if !gjson.ValidBytes(body) {
		conf.RecordFailure(net)
		a.versionQueue.In <- p
		log.Task.Warn("versionParse Error: invalid JSON response body")

		return
	}

	transfers := make([]transfer, 0)
	for _, trans := range gjson.ParseBytes(body).Array() {
		tsNano := trans.Get("timestamp").Int() * 1000
		timestamp := time.Unix(tsNano/1e9, tsNano%1e9)
		a.recordEarliestBlockTs(timestamp.Unix())

		ver := trans.Get("version").Int()
		hash := trans.Get("hash").String()
		addrOwner := make(map[string]string)                                         // [address] => owner address
		addrType := make(map[string]model.TradeType)                                 // [address] => tradeType
		amtAddrMap := map[string]map[aptAmount]string{"deposit": {}, "withdraw": {}} // [amount] => address
		aptEvents := make([]aptEvent, 0)
		trans.Get("changes").ForEach(func(_, v gjson.Result) bool {
			if v.Get("type").String() != "write_resource" {

				return true
			}

			data := v.Get("data")
			if data.Get("type").String() == "0x1::fungible_asset::FungibleStore" {
				addr := v.Get("address").String()
				switch data.Get("data.metadata.inner").String() {
				case conf.UsdtAptos:
					addrType[addr] = model.UsdtAptos
				case conf.UsdcAptos:
					addrType[addr] = model.UsdcAptos
				}
			}
			if data.Get("type").String() == "0x1::object::ObjectCore" {
				addrOwner[v.Get("address").String()] = data.Get("data.owner").String()
			}

			return true
		})
		trans.Get("events").ForEach(func(_, v gjson.Result) bool {
			amount := v.Get("data.amount").String()
			amt, err := decimal.NewFromString(amount)
			if err != nil {

				return true
			}

			address := v.Get("data.store").String()
			switch v.Get("type").String() {
			case "0x1::fungible_asset::Deposit":
				aptEvents = append(aptEvents, aptEvent{Amount: amt, Address: address, Action: "deposit"})
				amtAddrMap["deposit"][aptAmount{Amount: amount, Type: addrType[address]}] = address
			case "0x1::fungible_asset::Withdraw":
				amtAddrMap["withdraw"][aptAmount{Amount: amount, Type: addrType[address]}] = address
				aptEvents = append(aptEvents, aptEvent{Amount: amt, Address: address, Action: "withdraw"})
			}
			return true
		})

		// 针对 一个withdraw 对应 一个deposit 且数额相同的情况
		for amt, to := range amtAddrMap["deposit"] {
			from, ok := amtAddrMap["withdraw"][amt]
			if !ok {

				continue
			}

			amount, ok := new(big.Int).SetString(amt.Amount, 10)
			if !ok {

				continue
			}

			tradeType, ok := addrType[to]
			if !ok {

				continue
			}

			transfers = append(transfers, transfer{
				Network:     net,
				TxHash:      hash,
				Amount:      decimal.NewFromBigInt(amount, model.GetTradeDecimal(tradeType)),
				FromAddress: a.padAddressLeadingZeros(addrOwner[from]),
				RecvAddress: a.padAddressLeadingZeros(addrOwner[to]),
				Timestamp:   timestamp,
				TradeType:   tradeType,
				BlockNum:    ver,
			})
		}

		// 针对 一个withdraw 对应 多个deposit(数额累计等于 withdraw) 的情况
		processEvents := func(tradeType model.TradeType, events []aptEvent) ([]aptEvent, map[string]string) {
			deposits := make([]aptEvent, 0)
			withdraws := make(map[decimal.Decimal]aptEvent)
			fromMap := make(map[string]string)

			// 分类事件
			for _, e := range events {
				if addrType[e.Address] == tradeType {
					if e.Action == "deposit" {
						deposits = append(deposits, e)
					}
					if e.Action == "withdraw" {
						withdraws[e.Amount] = e
					}
				}
			}

			// 穷举计算匹配关系，只穷举 A + B = C 的情况，实际上还存在 A + B + C + ... = D
			// 大部分这种情况都是合约 swap 等交易，非普通人1对1转账，所以选择忽视
			for k1, e1 := range deposits {
				for k2, e2 := range deposits {
					if k1 == k2 {
						continue
					}
					for sum, e3 := range withdraws {
						if e1.Amount.Add(e2.Amount).Equal(sum) {
							fromMap[e1.Address] = e3.Address
						}
					}
				}
			}

			return deposits, fromMap
		}
		generateTransfers := func(deposits []aptEvent, fromMap map[string]string, t model.TradeType, decimals int32) {
			for _, to := range deposits {
				if from, ok := fromMap[to.Address]; ok {
					transfers = append(transfers, transfer{
						Network:     net,
						TxHash:      hash,
						Amount:      decimal.NewFromBigInt(to.Amount.BigInt(), decimals),
						FromAddress: a.padAddressLeadingZeros(addrOwner[from]),
						RecvAddress: a.padAddressLeadingZeros(addrOwner[to.Address]),
						Timestamp:   timestamp,
						TradeType:   t,
						BlockNum:    ver,
					})
				}
			}
		}

		// 处理 USDT
		usdtDeposits, usdtFrom := processEvents(model.UsdtAptos, aptEvents)
		generateTransfers(usdtDeposits, usdtFrom, model.UsdtAptos, model.GetTradeDecimal(model.UsdtAptos))

		// 处理 USDC
		usdcDeposits, usdcFrom := processEvents(model.UsdcAptos, aptEvents)
		generateTransfers(usdcDeposits, usdcFrom, model.UsdcAptos, model.GetTradeDecimal(model.UsdcAptos))
	}

	if len(transfers) > 0 {

		transferQueue.In <- transfers
	}

	log.Task.Info(fmt.Sprintf("区块扫描完成(Aptos) %d.%d 成功率：%s", p.Start, p.Limit, conf.GetSuccessRate(net)))
}

func (a *aptos) padAddressLeadingZeros(addr string) string {
	addr = strings.TrimPrefix(addr, "0x")
	addr = strings.Repeat("0", 64-len(addr)) + addr

	return "0x" + addr
}

func (a *aptos) tradeConfirmHandle(ctx context.Context) {
	var orders = getConfirmingOrders(model.GetNetworkTrades(conf.Aptos))
	var wg sync.WaitGroup

	var handle = func(o model.Order) {
		req, _ := http.NewRequestWithContext(ctx, "GET", model.Endpoint(conf.Aptos)+"v1/transactions/by_hash/"+o.RefHash, nil)
		resp, err := a.client.Do(req)
		if err != nil {
			log.Task.Warn("aptos tradeConfirmHandle Error sending request:", err)

			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Task.Warn("aptos tradeConfirmHandle Error response status code:", resp.StatusCode)

			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Task.Warn("aptos tradeConfirmHandle Error reading response body:", err)

			return
		}

		data := gjson.ParseBytes(body)
		if data.Get("error_code").Exists() {
			log.Task.Warn("aptos tradeConfirmHandle Error:", data.Get("message").String())

			return
		}

		if data.Get("version").String() != "" &&
			data.Get("success").Bool() &&
			data.Get("vm_status").String() == "Executed successfully" {

			markFinalConfirmed(o)
		}
	}

	for _, order := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()

			handle(order)
		}()
	}

	wg.Wait()
}

func (a *aptos) recordEarliestBlockTs(unix int64) {
	for {
		current := a.earliestBlockTs.Load()
		if current != 0 && unix >= current {
			return
		}
		if a.earliestBlockTs.CompareAndSwap(current, unix) {
			return
		}
	}
}
