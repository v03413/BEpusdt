package task

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

type ton struct {
	client         *http.Client
	ltOffset       int
	lastLt         atomic.Int64
	apiKeyCursor   atomic.Uint64
	apiKeyCooldown sync.Map
	seen           sync.Map
}

var tn ton

func tonInit() {
	tn = newTon()
	Register(Task{Callback: tn.syncEvents, Duration: time.Second * 5})
	Register(Task{Callback: tn.tradeConfirmHandle, Duration: time.Second * 5})
}

func newTon() ton {
	return ton{
		client:   utils.NewHttpClient(),
		ltOffset: 20,
	}
}

// syncEvents 同步 TON 活跃收款地址的原生币交易和 Jetton 转账。
func (t *ton) syncEvents(ctx context.Context) {
	if syncBreak(conf.Ton, 0) {
		return
	}

	addresses, startDate := t.getActiveAddresses()
	if len(addresses) == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 3)
	for _, address := range addresses {
		wg.Add(1)
		sem <- struct{}{}
		go func(address string) {
			defer wg.Done()
			defer func() { <-sem }()

			t.fetchAccountTransfers(ctx, address, startDate)
		}(address)
	}

	wg.Wait()
}

// getActiveAddresses 获取当前需要扫描的 TON 地址和最早订单时间。
func (t *ton) getActiveAddresses() ([]string, int64) {
	trades := model.GetNetworkTrades(conf.Ton)
	data := make(map[string]struct{})
	var startDate int64

	var orders []model.Order
	model.Db.Model(&model.Order{}).
		Where("status = ? and trade_type in (?)", model.OrderStatusWaiting, trades).
		Find(&orders)
	for _, order := range orders {
		address := model.NormalizeTradeAddress(order.TradeType, order.Address)
		if address != "" {
			data[address] = struct{}{}
		}
		created := order.CreatedAt.Time().Unix()
		if startDate == 0 || created < startDate {
			startDate = created
		}
	}

	var wallets []model.Wallet
	model.Db.Model(&model.Wallet{}).
		Where("status = ? and trade_type in (?)", model.WaStatusEnable, trades).
		Find(&wallets)
	for _, wallet := range wallets {
		address := model.NormalizeTradeAddress(model.TradeType(wallet.TradeType), wallet.MatchAddr)
		if address == "" {
			address = model.NormalizeTradeAddress(model.TradeType(wallet.TradeType), wallet.Address)
		}
		if address != "" {
			data[address] = struct{}{}
		}
	}

	addresses := make([]string, 0, len(data))
	for address := range data {
		addresses = append(addresses, address)
	}

	return addresses, startDate
}

// fetchAccountTransfers 从 TON Center V3 拉取单个账户的原生 TON 和 Jetton 入账。
func (t *ton) fetchAccountTransfers(ctx context.Context, address string, startDate int64) {
	transfers := make([]transfer, 0)
	if nativeTransfers, ok := t.fetchNativeTransfers(ctx, address, startDate); ok {
		transfers = append(transfers, nativeTransfers...)
	}
	for _, contract := range t.tonJettonContracts() {
		jettonTransfers, ok := t.fetchJettonTransfers(ctx, address, contract, startDate)
		if ok {
			transfers = append(transfers, jettonTransfers...)
		}
	}

	if len(transfers) > 0 {
		transferQueue.In <- transfers
	}

	log.Task.Info(fmt.Sprintf("区块扫描完成(Ton) %s 成功率：%s", model.FormatNetworkAddress(conf.Ton, address), conf.GetSuccessRate(conf.Ton)))
}

// fetchNativeTransfers 从 TON Center V3 查询账户原生 TON 入账交易。
func (t *ton) fetchNativeTransfers(ctx context.Context, address string, startDate int64) ([]transfer, bool) {
	query := url.Values{}
	query.Set("account", address)
	query.Set("limit", "100")
	query.Set("sort", "desc")
	if startDate > 0 {
		query.Set("start_utime", fmt.Sprintf("%d", startDate-60))
	}

	body, ok := t.get(ctx, "/transactions", query)
	if !ok {
		return nil, false
	}

	return t.parseNativeTransfers(body, address), true
}

// fetchJettonTransfers 从 TON Center V3 查询账户指定 Jetton 入账交易。
func (t *ton) fetchJettonTransfers(ctx context.Context, address string, contract string, startDate int64) ([]transfer, bool) {
	query := url.Values{}
	query.Set("owner_address", address)
	query.Set("jetton_master", contract)
	query.Set("direction", "in")
	query.Set("limit", "100")
	query.Set("sort", "desc")
	if startDate > 0 {
		query.Set("start_utime", fmt.Sprintf("%d", startDate-60))
	}

	body, ok := t.get(ctx, "/jetton/transfers", query)
	if !ok {
		return nil, false
	}

	return t.parseJettonTransfers(body), true
}

// parseNativeTransfers 解析 TON Center V3 账户交易为原生 TON 转账。
func (t *ton) parseNativeTransfers(body []byte, account string) []transfer {
	transfers := make([]transfer, 0)
	account = t.normalizeAddress(account)

	for _, tx := range gjson.GetBytes(body, "transactions").Array() {
		if tx.Get("description.aborted").Bool() {
			continue
		}

		msg := tx.Get("in_msg")
		recv := t.normalizeAddress(msg.Get("destination").String())
		if recv == "" || recv != account {
			continue
		}

		amount, ok := new(big.Int).SetString(msg.Get("value").String(), 10)
		if !ok || amount.Sign() <= 0 {
			continue
		}

		lt := t.parseLt(tx.Get("lt").String())
		t.recordLastLt(lt)

		tr := transfer{
			Network:     conf.Ton,
			TxHash:      tx.Get("hash").String(),
			Amount:      decimal.NewFromBigInt(amount, model.GetTradeDecimal(model.TonTon)),
			FromAddress: t.normalizeAddress(msg.Get("source").String()),
			RecvAddress: recv,
			Timestamp:   time.Unix(tx.Get("now").Int(), 0),
			TradeType:   model.TonTon,
			BlockNum:    int(lt),
		}
		if t.validTransfer(tr) {
			transfers = append(transfers, tr)
		}
	}

	return transfers
}

// parseJettonTransfers 解析 TON Center V3 Jetton 转账为统一转账结构。
func (t *ton) parseJettonTransfers(body []byte) []transfer {
	transfers := make([]transfer, 0)

	for _, item := range gjson.GetBytes(body, "jetton_transfers").Array() {
		if item.Get("transaction_aborted").Bool() {
			continue
		}

		contract := t.normalizeAddress(item.Get("jetton_master").String())
		tradeType, ok := model.GetContractTrade(contract)
		if !ok {
			continue
		}

		amount, ok := new(big.Int).SetString(item.Get("amount").String(), 10)
		if !ok || amount.Sign() <= 0 {
			continue
		}

		lt := t.parseLt(item.Get("transaction_lt").String())
		t.recordLastLt(lt)

		tr := transfer{
			Network:     conf.Ton,
			TxHash:      item.Get("transaction_hash").String(),
			Amount:      decimal.NewFromBigInt(amount, model.GetTradeDecimal(tradeType)),
			FromAddress: t.normalizeAddress(item.Get("source").String()),
			RecvAddress: t.normalizeAddress(item.Get("destination").String()),
			Timestamp:   time.Unix(item.Get("transaction_now").Int(), 0),
			TradeType:   tradeType,
			BlockNum:    int(lt),
		}
		if t.validTransfer(tr) {
			transfers = append(transfers, tr)
		}
	}

	return transfers
}

// validTransfer 过滤字段不完整或本轮已经处理过的 TON 转账。
func (t *ton) validTransfer(tr transfer) bool {
	if tr.TxHash == "" || tr.FromAddress == "" || tr.RecvAddress == "" || tr.Amount.IsZero() {
		return false
	}

	key := strings.Join([]string{tr.TxHash, string(tr.TradeType), tr.FromAddress, tr.RecvAddress, tr.Amount.String()}, ":")
	if _, loaded := t.seen.LoadOrStore(key, struct{}{}); loaded {
		return false
	}

	return true
}

// tradeConfirmHandle 确认 TON 网络中处于确认中的订单。
func (t *ton) tradeConfirmHandle(ctx context.Context) {
	var orders = getConfirmingOrders(model.GetNetworkTrades(conf.Ton))
	var wg sync.WaitGroup

	for _, order := range orders {
		wg.Add(1)
		go func(o model.Order) {
			defer wg.Done()

			if model.GetC(model.BlockOffsetConfirm) == "1" {
				last := t.lastLt.Load()
				if last == 0 || int(last)-o.RefBlockNum < t.ltOffset {
					return
				}
			}

			t.confirmOrder(ctx, o)
		}(order)
	}

	wg.Wait()
}

// confirmOrder 从 TON Center V3 重新读取交易并确认订单最终状态。
func (t *ton) confirmOrder(ctx context.Context, order model.Order) {
	if order.TradeType == model.TonTon {
		t.confirmNativeOrder(ctx, order)
		return
	}

	t.confirmJettonOrder(ctx, order)
}

// confirmNativeOrder 通过交易哈希确认原生 TON 订单。
func (t *ton) confirmNativeOrder(ctx context.Context, order model.Order) {
	query := url.Values{}
	query.Set("account", model.NormalizeTradeAddress(order.TradeType, order.Address))
	query.Set("hash", order.RefHash)
	query.Set("limit", "1")

	body, ok := t.get(ctx, "/transactions", query)
	if !ok {
		return
	}

	for _, tx := range gjson.GetBytes(body, "transactions").Array() {
		if tx.Get("hash").String() != order.RefHash || tx.Get("description.aborted").Bool() {
			continue
		}

		msg := tx.Get("in_msg")
		if t.normalizeAddress(msg.Get("destination").String()) == model.NormalizeTradeAddress(order.TradeType, order.Address) {
			markFinalConfirmed(order)

			return
		}
	}
}

// confirmJettonOrder 通过交易 LT 和哈希确认 TON Jetton 订单。
func (t *ton) confirmJettonOrder(ctx context.Context, order model.Order) {
	confMap := model.GetAllTradeConfig()
	tradeConf, ok := confMap[string(order.TradeType)]
	if !ok || tradeConf.Contract == "" {
		return
	}

	contract := t.normalizeAddress(tradeConf.Contract)
	query := url.Values{}
	query.Set("owner_address", model.NormalizeTradeAddress(order.TradeType, order.Address))
	query.Set("jetton_master", contract)
	query.Set("direction", "in")
	query.Set("start_lt", fmt.Sprintf("%d", order.RefBlockNum))
	query.Set("end_lt", fmt.Sprintf("%d", order.RefBlockNum))
	query.Set("limit", "10")

	body, ok := t.get(ctx, "/jetton/transfers", query)
	if !ok {
		return
	}

	for _, item := range gjson.GetBytes(body, "jetton_transfers").Array() {
		if item.Get("transaction_hash").String() != order.RefHash || item.Get("transaction_aborted").Bool() {
			continue
		}

		if t.normalizeAddress(item.Get("destination").String()) == model.NormalizeTradeAddress(order.TradeType, order.Address) {
			markFinalConfirmed(order)

			return
		}
	}
}

// get 调用 TON Center V3 GET 接口并处理通用状态记录。
func (t *ton) get(ctx context.Context, path string, query url.Values) ([]byte, bool) {
	reqURL := t.endpoint(path)
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	allKeys := model.GetTonCenterV3ApiKeys()
	if len(allKeys) == 0 {
		body, statusCode, err := t.getOnce(ctx, reqURL, "")
		if err != nil {
			conf.RecordFailure(conf.Ton)
			log.Task.Warn("ton center Error sending request:", err)

			return nil, false
		}
		if statusCode != http.StatusOK {
			conf.RecordFailure(conf.Ton)
			log.Task.Warn("ton center Error response status code:", statusCode)

			return nil, false
		}

		conf.RecordSuccess(conf.Ton)

		return body, true
	}

	keys := t.availableApiKeys(allKeys)
	if len(keys) == 0 {
		conf.RecordFailure(conf.Ton)
		log.Task.Warn("ton center Error all Api Keys are cooling down")

		return nil, false
	}

	var lastStatusCode int
	for _, apiKey := range keys {
		body, statusCode, err := t.getOnce(ctx, reqURL, apiKey)
		if err != nil {
			conf.RecordFailure(conf.Ton)
			log.Task.Warn("ton center Error sending request:", err)

			return nil, false
		}

		lastStatusCode = statusCode
		if statusCode == http.StatusTooManyRequests {
			t.cooldownApiKey(apiKey)
			log.Task.Warn("ton center Api Key rate limited, switching key...")

			continue
		}
		if statusCode != http.StatusOK {
			conf.RecordFailure(conf.Ton)
			log.Task.Warn("ton center Error response status code:", statusCode)

			return nil, false
		}

		conf.RecordSuccess(conf.Ton)

		return body, true
	}

	conf.RecordFailure(conf.Ton)
	if lastStatusCode == http.StatusTooManyRequests {
		log.Task.Warn("ton center Error all Api Keys are rate limited")
	} else {
		log.Task.Warn("ton center Error no available Api Key")
	}

	return nil, false
}

// getOnce 使用指定 Api Key 对 TON Center V3 发起一次 GET 请求。
func (t *ton) getOnce(ctx context.Context, reqURL string, apiKey string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		log.Task.Warn("ton center Error creating request:", err)

		return nil, 0, err
	}
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

// availableApiKeys 返回当前未冷却的 TON Center V3 Api Key，并按轮换游标排序。
func (t *ton) availableApiKeys(keys []string) []string {
	now := time.Now()
	available := make([]string, 0, len(keys))
	for _, key := range keys {
		if until, ok := t.apiKeyCooldown.Load(key); ok {
			if cooldownUntil, ok := until.(time.Time); ok && now.Before(cooldownUntil) {
				continue
			}
			t.apiKeyCooldown.Delete(key)
		}
		available = append(available, key)
	}
	if len(available) == 0 {
		return nil
	}

	start := int(t.apiKeyCursor.Add(1) % uint64(len(available)))
	rotated := make([]string, 0, len(available))
	rotated = append(rotated, available[start:]...)
	rotated = append(rotated, available[:start]...)

	return rotated
}

// cooldownApiKey 将触发频率限制的 TON Center V3 Api Key 暂时冷却。
func (t *ton) cooldownApiKey(apiKey string) {
	t.apiKeyCooldown.Store(apiKey, time.Now().Add(10*time.Second))
}

// tonJettonContracts 获取 TON 网络已注册的 Jetton 主合约地址。
func (t *ton) tonJettonContracts() []string {
	confMap := model.GetAllTradeConfig()
	contracts := make([]string, 0)
	seen := make(map[string]struct{})
	for _, tradeConf := range confMap {
		if tradeConf.Network != conf.Ton || tradeConf.Native || tradeConf.Contract == "" {
			continue
		}

		contract := t.normalizeAddress(tradeConf.Contract)
		if contract == "" {
			continue
		}
		if _, ok := seen[contract]; ok {
			continue
		}

		seen[contract] = struct{}{}
		contracts = append(contracts, contract)
	}

	return contracts
}

// recordLastLt 记录当前进程已观察到的最大 TON 逻辑时间。
func (t *ton) recordLastLt(lt int64) {
	for {
		current := t.lastLt.Load()
		if lt <= current {
			return
		}
		if t.lastLt.CompareAndSwap(current, lt) {
			return
		}
	}
}

// parseLt 将 TON Center V3 的字符串 LT 转成整数。
func (t *ton) parseLt(value string) int64 {
	lt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return lt
}

// normalizeAddress 将 TON 地址统一转换为 raw 地址。
func (t *ton) normalizeAddress(address string) string {
	if normalized, ok := utils.NormalizeTonAddress(address); ok {
		return normalized
	}

	return strings.TrimSpace(address)
}

// endpoint 拼接当前配置的 TON Center V3 端点和请求路径。
func (t *ton) endpoint(path string) string {
	endpoint := strings.TrimRight(model.Endpoint(conf.Ton), "/")

	return endpoint + path
}
