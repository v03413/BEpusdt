package task

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"github.com/smallnest/chanx"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

// 参考文档
//  - https://solana.com/zh/docs/rpc
//  - https://github.com/solana-program/token/blob/6d18ff73b1dd30703a30b1ca941cb0f1d18c2b2a/program/src/instruction.rs

type solana struct {
	slotConfirmedOffset int
	lastSlotNum         int
	slotQueue           *chanx.UnboundedChan[int]
	client              *http.Client
}

type solanaTokenOwner struct {
	TradeType model.TradeType
	Address   string
}

var sol solana

func init() {
	sol = newSolana()
	Register(Task{Callback: sol.slotDispatch})
	Register(Task{Callback: sol.syncSlotForward, Duration: time.Second * 5})
	Register(Task{Callback: sol.tradeConfirmHandle, Duration: time.Second * 5})
	Register(Task{Callback: sol.reconcileWaitingOrders, Duration: time.Second * 15})
}

func newSolana() solana {
	return solana{
		slotConfirmedOffset: 60,
		lastSlotNum:         0,
		slotQueue:           chanx.NewUnboundedChan[int](context.Background(), 30),
		client:              utils.NewHttpClient(),
	}
}

func (s *solana) syncSlotForward(ctx context.Context) {
	if syncBreak(conf.Solana, s.slotQueue.Len()) {

		return
	}

	req, _ := http.NewRequestWithContext(ctx, "POST", model.Endpoint(conf.Solana), bytes.NewBuffer([]byte(`{"jsonrpc":"2.0","id":1,"method":"getSlot"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		log.Task.Warn("syncSlotForward Error sending request:", err)

		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Task.Warn("syncSlotForward Error response status code:", resp.StatusCode)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Task.Warn("syncSlotForward Error reading response body:", err)

		return
	}

	now := int(gjson.GetBytes(body, "result").Int())
	if now <= 0 {
		log.Task.Warn("syncSlotForward Error: invalid slot number:", now)

		return
	}

	if s.lastSlotNum == 0 { // 首次启动，往后追溯
		s.syncSlotBackward(now)
	}

	if now-s.lastSlotNum > cast.ToInt(model.GetC(model.BlockHeightMaxDiff)) { // 区块高度变化过大，强制丢块重扫
		s.lastSlotNum = now
	}

	if now == s.lastSlotNum { // 区块高度没有变化

		return
	}

	for n := s.lastSlotNum + 1; n <= now; n++ {
		// 待扫描区块入列

		s.slotQueue.In <- n
	}

	s.lastSlotNum = now
}

func (s *solana) syncSlotBackward(now int) {
	if now == 0 || s.lastSlotNum != 0 {

		return
	}

	var o model.Order
	trade := model.GetNetworkTrades(conf.Solana)
	model.Db.Model(&model.Order{}).Where("status = ? and trade_type in (?)", model.OrderStatusWaiting, trade).Order("created_at asc").Limit(1).Find(&o)
	if o.ID == 0 {

		return
	}

	// Solana 大概1秒3个区块（大概值，实际存在波动）
	num := int((time.Now().Unix() - o.CreatedAt.Time().Unix() + 1) * 3) // 计算需要反向扫描的区块数量

	go func() {
		ticker := time.NewTicker(time.Millisecond * 125)
		defer ticker.Stop()

		for i := 0; i < num; i++ {
			if syncBreak(conf.Solana, s.slotQueue.Len()) {

				return
			}

			s.slotQueue.In <- now - i

			<-ticker.C
		}
	}()
}

func (s *solana) slotDispatch(ctx context.Context) {
	p, err := ants.NewPoolWithFunc(3, s.slotParse)
	if err != nil {
		log.Task.Warn("Error creating pool:", err)

		return
	}

	defer p.Release()

	for {
		select {
		case slot := <-s.slotQueue.Out:
			if err := p.Invoke(slot); err != nil {
				s.slotQueue.In <- slot
				log.Task.Warn("slotDispatch Error invoking process slot:", err)
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				log.Task.Warn("slotDispatch context done:", err)
			}

			return
		}
	}
}

func (s *solana) slotParse(n any) {
	slot := n.(int)
	post := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"getBlock","params":[%d,{"encoding":"json","maxSupportedTransactionVersion":0,"transactionDetails":"full","rewards":false}]}`, slot))
	network := conf.Solana

	conf.RecordSuccess(network)
	resp, err := s.client.Post(model.Endpoint(conf.Solana), "application/json", bytes.NewBuffer(post))
	if err != nil {
		conf.RecordFailure(network)
		log.Task.Warn("slotParse Error sending request:", err)

		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		conf.RecordFailure(network)
		log.Task.Warn("slotParse Error response status code:", resp.StatusCode)

		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		conf.RecordFailure(network)
		s.slotQueue.In <- slot
		log.Task.Warn("slotParse Error reading response body:", err)

		return
	}

	timestamp := time.Unix(gjson.GetBytes(body, "result.blockTime").Int(), 0)

	for _, trans := range gjson.GetBytes(body, "result.transactions").Array() {
		hash := trans.Get("transaction.signatures.0").String()

		// 解析账号索引
		accountKeys := make([]string, 0)
		for _, key := range trans.Get("transaction.message.accountKeys").Array() {
			accountKeys = append(accountKeys, key.String())
		}
		for _, v := range []string{"readonly", "writable"} {
			for _, key := range trans.Get("meta.loadedAddresses." + v).Array() {
				accountKeys = append(accountKeys, key.String())
			}
		}

		// 查找SPL Token索引
		splTokenIndex := int64(-1)
		for i, v := range accountKeys {
			if v == conf.SolSplToken {
				splTokenIndex = int64(i)

				break
			}
		}

		// SPL Token的Mint地址，即不包含 Token 交易信息
		if splTokenIndex == -1 {

			continue
		}

		// 解析 Token 账户 【Token Wallet => Owner Wallet】
		tokenAccountMap := make(map[string]solanaTokenOwner)
		for _, v := range []string{"postTokenBalances", "preTokenBalances"} {
			for _, itm := range trans.Get("meta." + v).Array() {
				tradeType, ok := model.GetContractTrade(itm.Get("mint").String())
				if !ok || itm.Get("programId").String() != conf.SolSplToken {

					continue
				}

				tokenAccountMap[accountKeys[itm.Get("accountIndex").Int()]] = solanaTokenOwner{
					TradeType: tradeType,
					Address:   itm.Get("owner").String(),
				}
			}
		}

		transArr := make([]transfer, 0)

		// 解析外部指令
		for _, instr := range trans.Get("transaction.message.instructions").Array() {
			if instr.Get("programIdIndex").Int() != splTokenIndex {

				continue
			}

			transArr = append(transArr, s.parseTransfer(instr, accountKeys, tokenAccountMap))
		}

		// 解析内部指令
		for _, itm := range trans.Get("meta.innerInstructions").Array() {
			for _, instr := range itm.Get("instructions").Array() {
				if instr.Get("programIdIndex").Int() != splTokenIndex {

					continue
				}

				transArr = append(transArr, s.parseTransfer(instr, accountKeys, tokenAccountMap))
			}
		}

		// 过滤无关交易
		result := make([]transfer, 0)
		for _, t := range transArr {
			if t.FromAddress == "" || t.RecvAddress == "" || t.Amount.IsZero() {

				continue
			}

			t.TxHash = hash
			t.Network = conf.Solana
			t.BlockNum = slot
			t.Timestamp = timestamp

			result = append(result, t)
		}

		if len(result) > 0 {
			transferQueue.In <- result
		}
	}

	log.Task.Info(fmt.Sprintf("区块扫描完成(Solana) %d 成功率：%s", slot, conf.GetSuccessRate(network)))
}

func (s *solana) parseTransfer(instr gjson.Result, accountKeys []string, tokenAccountMap map[string]solanaTokenOwner) transfer {
	accounts := instr.Get("accounts").Array()
	trans := transfer{}
	if len(accounts) < 3 { // from to singer，至少存在3个账户索引，如果是多签则 > 3

		return trans
	}

	data := base58.Decode(instr.Get("data").String())
	dLen := len(data)
	if dLen < 9 {

		return trans
	}

	isTransfer := data[0] == 3 && dLen == 9
	isTransferChecked := data[0] == 12 && dLen == 10
	if !isTransfer && !isTransferChecked {

		return trans
	}

	var exp int32 = -6
	if isTransferChecked {
		exp = int32(data[9]) * -1
	}

	from, ok := tokenAccountMap[accountKeys[accounts[0].Int()]]
	if !ok {

		return trans
	}

	trans.FromAddress = from.Address
	trans.RecvAddress = tokenAccountMap[accountKeys[accounts[1].Int()]].Address
	if isTransferChecked {
		trans.RecvAddress = tokenAccountMap[accountKeys[accounts[2].Int()]].Address
	}

	buf := make([]byte, 8)
	copy(buf[:], data[1:9])
	number := binary.LittleEndian.Uint64(buf)
	b := new(big.Int)
	b.SetUint64(number)
	trans.TradeType = from.TradeType
	trans.Amount = decimal.NewFromBigInt(b, exp)

	return trans
}

func (s *solana) tradeConfirmHandle(ctx context.Context) {
	var orders = getConfirmingOrders(model.GetNetworkTrades(conf.Solana))
	var wg sync.WaitGroup

	var handle = func(o model.Order) {
		if model.GetC(model.BlockOffsetConfirm) == "1" {
			if s.lastSlotNum == 0 {
				return
			}
			if s.lastSlotNum-o.RefBlockNum < s.slotConfirmedOffset {
				return
			}
		}

		post := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"getSignatureStatuses","params":[["%s"],{"searchTransactionHistory":true}]}`, o.RefHash))
		req, _ := http.NewRequestWithContext(ctx, "POST", model.Endpoint(conf.Solana), bytes.NewBuffer(post))
		resp, err := s.client.Do(req)
		if err != nil {
			log.Task.Warn("solana tradeConfirmHandle Error sending request:", err)

			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Task.Warn("solana tradeConfirmHandle Error response status code:", resp.StatusCode)

			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Task.Warn("solana tradeConfirmHandle Error reading response body:", err)

			return
		}

		data := gjson.ParseBytes(body)
		if data.Get("error").Exists() {
			log.Task.Warn("solana tradeConfirmHandle Error:", data.Get("error").String())

			return
		}

		if data.Get("result.value.0.confirmationStatus").String() == "finalized" {

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

func (s *solana) reconcileWaitingOrders(ctx context.Context) {
	var orders []model.Order
	trades := model.GetNetworkTrades(conf.Solana)
	if len(trades) == 0 {
		return
	}

	status := []int{model.OrderStatusWaiting, model.OrderStatusExpired, model.OrderStatusFailed}
	model.Db.Where("status in (?) and trade_type in (?)", status, trades).
		Where("expired_at > ?", time.Now().Add(-24*time.Hour)).
		Find(&orders)

	if len(orders) == 0 {
		return
	}

	tokenAccounts := make(map[string][]string)
	for _, order := range orders {
		if ctx.Err() != nil {
			return
		}

		key := fmt.Sprintf("%s%s", order.Address, order.TradeType)
		accounts, ok := tokenAccounts[key]
		if !ok {
			var err error
			accounts, err = s.getTokenAccountsByOwner(ctx, order.Address, order.TradeType)
			if err != nil {
				log.Task.Warn("solana reconcile getTokenAccountsByOwner Error:", err)
				continue
			}
			tokenAccounts[key] = accounts
		}

		for _, account := range accounts {
			if s.reconcileOrderTokenAccount(ctx, order, account) {
				break
			}
		}
	}
}

func (s *solana) getTokenAccountsByOwner(ctx context.Context, owner string, tradeType model.TradeType) ([]string, error) {
	contract := ""
	if c, ok := model.GetAllTradeConfig()[string(tradeType)]; ok {
		contract = c.Contract
	}
	if contract == "" {
		return nil, nil
	}

	result, err := s.rpc(ctx, "getTokenAccountsByOwner", []any{
		owner,
		map[string]any{"mint": contract},
		map[string]any{"encoding": "jsonParsed"},
	})
	if err != nil {
		return nil, err
	}

	accounts := make([]string, 0)
	for _, item := range result.Get("value").Array() {
		pubkey := item.Get("pubkey").String()
		if pubkey != "" {
			accounts = append(accounts, pubkey)
		}
	}

	return accounts, nil
}

func (s *solana) reconcileOrderTokenAccount(ctx context.Context, order model.Order, account string) bool {
	result, err := s.rpc(ctx, "getSignaturesForAddress", []any{
		account,
		map[string]any{"limit": 50},
	})
	if err != nil {
		log.Task.Warn("solana reconcile getSignaturesForAddress Error:", err)
		return false
	}

	for _, sig := range result.Array() {
		if sig.Get("err").Exists() && sig.Get("err").Raw != "null" {
			continue
		}

		blockTime := sig.Get("blockTime").Int()
		if blockTime > 0 {
			ts := time.Unix(blockTime, 0)
			if !order.CreatedAt.Before(ts) || !order.ExpiredAt.After(ts) {
				continue
			}
		}

		hash := sig.Get("signature").String()
		if hash == "" {
			continue
		}

		result, err := s.rpc(ctx, "getTransaction", []any{
			hash,
			map[string]any{"encoding": "jsonParsed", "maxSupportedTransactionVersion": 0},
		})
		if err != nil {
			log.Task.Warn("solana reconcile getTransaction Error:", err)
			continue
		}
		if result.Get("meta.err").Exists() && result.Get("meta.err").Raw != "null" {
			continue
		}

		for _, t := range parseSolanaParsedTransfers(result) {
			t.TxHash = hash
			if t.BlockNum == 0 {
				t.BlockNum = int(result.Get("slot").Int())
			}
			if t.Timestamp.IsZero() {
				t.Timestamp = time.Unix(result.Get("blockTime").Int(), 0)
			}
			if !orderTransferMatch(order, t) {
				continue
			}

			if err := order.MarkConfirming(t.BlockNum, t.FromAddress, t.TxHash, t.Timestamp, t.Amount); err != nil {
				log.Task.Warn("solana reconcile mark order confirming failed:", err)
				return false
			}

			log.Task.Info(fmt.Sprintf("Solana 订单回查确认成功：%s %s", order.TradeId, t.TxHash))
			return true
		}
	}

	return false
}

func (s *solana) rpc(ctx context.Context, method string, params any) (gjson.Result, error) {
	post, err := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return gjson.Result{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", model.Endpoint(conf.Solana), bytes.NewBuffer(post))
	if err != nil {
		return gjson.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return gjson.Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return gjson.Result{}, fmt.Errorf("response status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return gjson.Result{}, err
	}

	data := gjson.ParseBytes(body)
	if data.Get("error").Exists() {
		return gjson.Result{}, fmt.Errorf("%s", data.Get("error").String())
	}

	return data.Get("result"), nil
}

func parseSolanaParsedTransfers(tx gjson.Result) []transfer {
	tokenAccountMap := make(map[string]solanaTokenOwner)
	for _, v := range []string{"postTokenBalances", "preTokenBalances"} {
		for _, item := range tx.Get("meta." + v).Array() {
			tradeType, ok := model.GetContractTrade(item.Get("mint").String())
			if !ok || item.Get("programId").String() != conf.SolSplToken {
				continue
			}

			tokenAccountMap[item.Get("accountIndex").String()] = solanaTokenOwner{
				TradeType: tradeType,
				Address:   item.Get("owner").String(),
			}
		}
	}

	transfers := make([]transfer, 0)
	instructions := tx.Get("transaction.message.instructions").Array()
	for _, inner := range tx.Get("meta.innerInstructions").Array() {
		instructions = append(instructions, inner.Get("instructions").Array()...)
	}

	for _, instr := range instructions {
		if instr.Get("programId").String() != conf.SolSplToken {
			continue
		}

		parsed := instr.Get("parsed")
		if !parsed.Exists() {
			continue
		}

		typ := parsed.Get("type").String()
		if typ != "transfer" && typ != "transferChecked" {
			continue
		}

		info := parsed.Get("info")
		source := info.Get("source").String()
		destination := info.Get("destination").String()
		from, ok := tokenAccountOwnerByPubkey(tx, tokenAccountMap, source)
		if !ok {
			continue
		}
		to, ok := tokenAccountOwnerByPubkey(tx, tokenAccountMap, destination)
		if !ok {
			continue
		}

		amountRaw := info.Get("tokenAmount.amount").String()
		decimals := info.Get("tokenAmount.decimals").Int()
		if amountRaw == "" {
			amountRaw = info.Get("amount").String()
			decimals = 6
		}

		amountInt, ok := new(big.Int).SetString(amountRaw, 10)
		if !ok {
			continue
		}

		transfers = append(transfers, transfer{
			Network:     conf.Solana,
			Amount:      decimal.NewFromBigInt(amountInt, -int32(decimals)),
			FromAddress: from.Address,
			RecvAddress: to.Address,
			Timestamp:   time.Unix(tx.Get("blockTime").Int(), 0),
			TradeType:   from.TradeType,
			BlockNum:    int(tx.Get("slot").Int()),
		})
	}

	return transfers
}

func tokenAccountOwnerByPubkey(tx gjson.Result, tokenAccountMap map[string]solanaTokenOwner, pubkey string) (solanaTokenOwner, bool) {
	accountKeys := tx.Get("transaction.message.accountKeys").Array()
	for i, key := range accountKeys {
		if key.Get("pubkey").String() == pubkey || key.String() == pubkey {
			owner, ok := tokenAccountMap[fmt.Sprintf("%d", i)]
			return owner, ok
		}
	}

	return solanaTokenOwner{}, false
}
