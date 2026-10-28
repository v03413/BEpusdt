package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/model"
)

// 按收款代币账户检测的测试：一个 httptest 服务模拟两个节点，/main 是主节点，/poll 是签名轮询节点。
// 交易取自主网 fixture：success_v0_inner_transfer.json 是 GMCJvYGf… 付给 watchTestOwner 390.990618 USDC，
// success_legacy_transfer_checked.json 是付给别人的 USDC，用作收款账户上的无关交易。

const (
	watchTestOwner  = "F7p3dFrjRTbtRp8FRF6qHLomXbKRBzpvBLjtQcfcgmNe"
	watchTestAmount = "390.990618"
	watchTestPay    = "success_v0_inner_transfer.json"
	watchTestOther  = "success_legacy_transfer_checked.json"
)

type watchTestSig struct {
	sig       string
	slot      int64
	blockTime int64
	err       string // 签名列表里的 err，默认 null
	fixture   string // getTransaction 返回的交易；为空时返回 null（节点还没有这笔交易）
	failMeta  bool   // 交易的 meta.err 不为空
}

type watchTestRPC struct {
	t      *testing.T
	mu     sync.Mutex
	lists  map[string][]watchTestSig // 代币账户 -> 签名，从新到旧
	extra  map[string]string         // 钱包 -> getTokenAccountsByOwner 额外返回的代币账户
	status map[string]int            // 节点路径 -> 强制返回的 HTTP 状态码
	calls  map[string]int            // "节点路径 方法" 和 "getTransaction 签名" 的调用次数
}

func newWatchTestRPC(t *testing.T) *watchTestRPC {
	return &watchTestRPC{
		t:      t,
		lists:  make(map[string][]watchTestSig),
		extra:  make(map[string]string),
		status: make(map[string]int),
		calls:  make(map[string]int),
	}
}

func (r *watchTestRPC) set(account string, sigs ...watchTestSig) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lists[account] = sigs
}

func (r *watchTestRPC) setStatus(path string, code int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status[path] = code
}

func (r *watchTestRPC) count(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.calls[key]
}

func (r *watchTestRPC) find(sig string) (watchTestSig, bool) {
	for _, list := range r.lists {
		for _, itm := range list {
			if itm.sig == sig {
				return itm, true
			}
		}
	}

	return watchTestSig{}, false
}

func (r *watchTestRPC) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/notify" {
		_, _ = w.Write([]byte("ok"))

		return
	}

	body, _ := io.ReadAll(req.Body)
	method := gjson.GetBytes(body, "method").String()
	param := gjson.GetBytes(body, "params.0").String()

	r.mu.Lock()
	r.calls[req.URL.Path+" "+method]++
	if method == "getTransaction" {
		r.calls["getTransaction "+param]++
	}
	status := r.status[req.URL.Path]
	r.mu.Unlock()
	if status != 0 {
		w.WriteHeader(status)

		return
	}

	switch method {
	case "getSignaturesForAddress":
		r.mu.Lock()
		list := r.lists[param]
		r.mu.Unlock()

		start := 0
		if before := gjson.GetBytes(body, "params.1.before").String(); before != "" {
			start = len(list)
			for i, itm := range list {
				if itm.sig == before {
					start = i + 1
				}
			}
		}
		end := min(start+int(gjson.GetBytes(body, "params.1.limit").Int()), len(list))

		items := make([]string, 0)
		for _, itm := range list[start:end] {
			errJSON := itm.err
			if errJSON == "" {
				errJSON = "null"
			}
			items = append(items, fmt.Sprintf(`{"signature":"%s","slot":%d,"err":%s,"memo":null,"blockTime":%d,"confirmationStatus":"confirmed"}`, itm.sig, itm.slot, errJSON, itm.blockTime))
		}
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","result":[%s],"id":1}`, strings.Join(items, ","))
	case "getTransaction":
		r.mu.Lock()
		itm, ok := r.find(param)
		r.mu.Unlock()
		if !ok || itm.fixture == "" {
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","result":null,"id":1}`))

			return
		}
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","result":%s,"id":1}`, watchTestTransaction(r.t, itm))
	case "getSignatureStatuses":
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","result":{"context":{"slot":%d},"value":[{"confirmationStatus":"finalized","confirmations":null,"err":null,"slot":1,"status":{"Ok":null}}]},"id":1}`, int64(1)<<40)
	case "getTokenAccountsByOwner":
		if req.URL.Path == "/poll" { // 模拟 publicnode 免费节点不提供该方法
			w.WriteHeader(http.StatusForbidden)

			return
		}
		ata, _ := solanaAssociatedTokenAccount(param, gjson.GetBytes(body, "params.1.mint").String())
		value := fmt.Sprintf(`{"pubkey":"%s","account":{}}`, ata)
		r.mu.Lock()
		if extra := r.extra[param]; extra != "" {
			value += fmt.Sprintf(`,{"pubkey":"%s","account":{}}`, extra)
		}
		r.mu.Unlock()
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","result":{"context":{"slot":1},"value":[%s]},"id":1}`, value)
	default:
		r.t.Errorf("unexpected RPC call %s on %s", method, req.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	}
}

// watchTestTransaction 把 fixture 里的交易改成指定签名、slot 和时间，按 getTransaction 的格式返回
func watchTestTransaction(t *testing.T, itm watchTestSig) string {
	t.Helper()

	dec := json.NewDecoder(bytes.NewReader([]byte(gjson.GetBytes(readSolanaFixture(t, itm.fixture), "result.transactions.0").Raw)))
	dec.UseNumber()
	var trans map[string]any
	if err := dec.Decode(&trans); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	trans["transaction"].(map[string]any)["signatures"].([]any)[0] = itm.sig
	if itm.failMeta {
		trans["meta"].(map[string]any)["err"] = map[string]any{"InstructionError": []any{1, map[string]any{"Custom": 1}}}
	}
	trans["slot"] = itm.slot
	trans["blockTime"] = itm.blockTime

	out, err := json.Marshal(trans)
	if err != nil {
		t.Fatalf("encode transaction: %v", err)
	}

	return string(out)
}

func initWatchTest(t *testing.T, rpc *watchTestRPC) (*solana, string) {
	t.Helper()

	srv := httptest.NewServer(rpc)
	t.Cleanup(srv.Close)
	initTaskTestEnv(t, map[model.ConfKey]string{
		model.RpcEndpointSolana:     srv.URL + "/main",
		model.RpcEndpointSolanaPoll: srv.URL + "/poll",
		model.SolanaScanMode:        solanaScanModeAddress,
	})
	solanaCooldown.Clear()
	s := newSolana()

	return &s, srv.URL
}

func watchTestATA(t *testing.T) string {
	t.Helper()

	ata, err := solanaAssociatedTokenAccount(watchTestOwner, model.GetTradeContract(model.UsdcSolana))
	if err != nil {
		t.Fatalf("derive ATA: %v", err)
	}

	return ata
}

func createWatchTestOrder(t *testing.T, srvURL, tradeID string, status int, created, expired time.Time) model.Order {
	t.Helper()

	return createTestOrder(t, model.Order{
		OrderId:   "merchant-" + tradeID,
		TradeId:   tradeID,
		TradeType: model.UsdcSolana,
		Rate:      "7.00",
		Amount:    watchTestAmount,
		Money:     "2736.93",
		Address:   watchTestOwner,
		Status:    status,
		ApiType:   model.OrderApiTypeEpusdt,
		NotifyUrl: srvURL + "/notify",
		ExpiredAt: expired,
	}, created)
}

// createWatchTestIdleOrder 建一个金额不同、不会被测试交易匹配的等待支付订单，让收款账户在第一笔付款入账后仍被检测
func createWatchTestIdleOrder(t *testing.T, srvURL string, now time.Time) {
	t.Helper()

	o := createWatchTestOrder(t, srvURL, "trade-idle", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	model.Db.Model(&o).Update("amount", "1.000001")
}

func watchTestPayment(sig string, slot int64, at time.Time) watchTestSig {
	return watchTestSig{sig: sig, slot: slot, blockTime: at.Unix(), fixture: watchTestPay}
}

func decimalOf(t *testing.T, v string) decimal.Decimal {
	t.Helper()

	d, err := decimal.NewFromString(v)
	if err != nil {
		t.Fatalf("decimal %q: %v", v, err)
	}

	return d
}

func getWatchTestCursor(t *testing.T, account string) model.SolanaCursor {
	t.Helper()

	c, err := model.GetSolanaCursor(account)
	if err != nil {
		t.Fatalf("get cursor: %v", err)
	}

	return c
}

func TestSolanaWatchDetectsPaymentAndConfirms(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	ata := watchTestATA(t)
	rpc.set(ata, watchTestPayment("pay-1", 5000, now.Add(-30*time.Second)))

	s.watchPayments(t.Context())

	got := getTestOrder(t, order.TradeId)
	if got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" || got.RefBlockNum != 5000 ||
		got.FromAddress != "GMCJvYGf5Ex2ARiMquaBDqU6iKM8uiEQkB8jCnoNfHpC" {
		t.Fatalf("order = status %d ref %s slot %d from %s, want confirming with pay-1", got.Status, got.RefHash, got.RefBlockNum, got.FromAddress)
	}
	if got.ConfirmedAt == nil || got.ConfirmedAt.Unix() != now.Add(-30*time.Second).Unix() {
		t.Fatalf("confirmed_at = %v, want the transaction time", got.ConfirmedAt)
	}
	if c := getWatchTestCursor(t, ata); c.Slot != 5000-solanaWatchSlotOverlap || c.Signature != "pay-1" || c.Owner != watchTestOwner {
		t.Fatalf("cursor = %+v, want slot %d after pay-1", c, 5000-solanaWatchSlotOverlap)
	}

	// 签名走轮询节点，交易和代币账户列表走主节点
	if rpc.count("/poll getSignaturesForAddress") == 0 || rpc.count("/main getSignaturesForAddress") != 0 {
		t.Fatalf("getSignaturesForAddress poll=%d main=%d, want poll only", rpc.count("/poll getSignaturesForAddress"), rpc.count("/main getSignaturesForAddress"))
	}
	if rpc.count("/main getTransaction") != 1 || rpc.count("/poll getTransaction") != 0 {
		t.Fatalf("getTransaction main=%d poll=%d, want main once", rpc.count("/main getTransaction"), rpc.count("/poll getTransaction"))
	}

	// 下一轮不再下载同一笔交易
	s.watchPayments(t.Context())
	if n := rpc.count("getTransaction pay-1"); n != 1 {
		t.Fatalf("getTransaction pay-1 = %d, want 1", n)
	}

	s.tradeConfirmHandle(t.Context())
	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusSuccess {
		t.Fatalf("order status = %d after confirm, want success", got.Status)
	}
	waitOrderNotified(t, order.TradeId)
}

func TestSolanaWatchSkipsFailedTransactions(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	ata := watchTestATA(t)

	failedInList := watchTestPayment("failed-in-list", 5001, now.Add(-20*time.Second))
	failedInList.err = `{"InstructionError":[1,{"Custom":6001}]}`
	failedInMeta := watchTestPayment("failed-in-meta", 5000, now.Add(-30*time.Second))
	failedInMeta.failMeta = true
	rpc.set(ata, failedInList, failedInMeta)

	s.watchPayments(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusWaiting || got.RefHash != "" {
		t.Fatalf("order = status %d ref %q, want waiting", got.Status, got.RefHash)
	}
	if n := rpc.count("getTransaction failed-in-list"); n != 0 {
		t.Fatalf("downloaded a transaction marked failed in the signature list %d times", n)
	}
	if n := rpc.count("getTransaction failed-in-meta"); n != 1 {
		t.Fatalf("getTransaction failed-in-meta = %d, want 1", n)
	}
	if c := getWatchTestCursor(t, ata); c.Slot != 5001-solanaWatchSlotOverlap {
		t.Fatalf("cursor slot = %d, want both failed signatures done", c.Slot)
	}
}

func TestSolanaWatchStuckTransactionDoesNotBlockNewPayment(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-2*time.Minute), now.Add(8*time.Minute))
	ata := watchTestATA(t)

	createWatchTestIdleOrder(t, srvURL, now)
	stuck := watchTestSig{sig: "stuck", slot: 4000, blockTime: now.Add(-90 * time.Second).Unix()} // 两个节点都还查不到
	rpc.set(ata, watchTestPayment("pay-1", 6000, now.Add(-10*time.Second)), stuck)

	s.watchPayments(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" {
		t.Fatalf("order = status %d ref %q, want confirming with pay-1 while another transaction is stuck", got.Status, got.RefHash)
	}
	if rpc.count("/main getTransaction") != 2 || rpc.count("/poll getTransaction") != 1 {
		t.Fatalf("getTransaction main=%d poll=%d, want the null result retried on the poll node", rpc.count("/main getTransaction"), rpc.count("/poll getTransaction"))
	}
	if c := getWatchTestCursor(t, ata); c.Slot != 4000 {
		t.Fatalf("cursor slot = %d, want it held at the stuck transaction 4000", c.Slot)
	}

	// 退避期间不再请求这笔交易
	s.watchPayments(t.Context())
	if n := rpc.count("getTransaction stuck"); n != 2 {
		t.Fatalf("getTransaction stuck = %d during backoff, want 2", n)
	}

	// 节点有了这笔交易、退避到期后补上，低水位前进
	stuck.fixture = watchTestOther
	rpc.set(ata, watchTestPayment("pay-1", 6000, now.Add(-10*time.Second)), stuck)
	s.watch.mu.Lock()
	s.watch.signatures["stuck"].retryAt = time.Time{}
	s.watch.mu.Unlock()
	s.watchPayments(t.Context())
	if c := getWatchTestCursor(t, ata); c.Slot != 6000-solanaWatchSlotOverlap {
		t.Fatalf("cursor slot = %d after the stuck transaction was processed, want %d", c.Slot, 6000-solanaWatchSlotOverlap)
	}
}

func TestSolanaWatchGivesUpOnTransactionAfterLimit(t *testing.T) {
	w := newSolanaWatch()
	start := time.Now()
	w.retryLater("sig", start, fmt.Errorf("not found"))
	if w.isDone("sig") || w.due("sig", start) {
		t.Fatalf("signature should wait for its backoff")
	}

	w.retryLater("sig", start.Add(solanaWatchGiveUp), fmt.Errorf("not found"))
	if !w.isDone("sig") {
		t.Fatalf("signature should be given up after %s", solanaWatchGiveUp)
	}
}

func TestSolanaWatchResumesFromCursorAfterRestart(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-3*time.Minute), now.Add(7*time.Minute))
	ata := watchTestATA(t)

	unrelated := watchTestSig{sig: "unrelated", slot: 1800, blockTime: now.Add(-2 * time.Minute).Unix(), fixture: watchTestOther}
	stuck := watchTestSig{sig: "pay-1", slot: 2000, blockTime: now.Add(-time.Minute).Unix()}
	rpc.set(ata, stuck, unrelated)

	s.watchPayments(t.Context())
	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusWaiting {
		t.Fatalf("order status = %d, want waiting before the transaction is available", got.Status)
	}
	if c := getWatchTestCursor(t, ata); c.Slot != 2000-solanaWatchSlotOverlap {
		t.Fatalf("cursor slot = %d, want %d", c.Slot, 2000-solanaWatchSlotOverlap)
	}

	// 重启：内存状态全部丢失，只剩库里的低水位
	rpc.set(ata, watchTestPayment("pay-1", 2000, now.Add(-time.Minute)), unrelated)
	restarted := newSolana()
	restarted.watchPayments(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" {
		t.Fatalf("order = status %d ref %q after restart, want confirming with pay-1", got.Status, got.RefHash)
	}
	if n := rpc.count("getTransaction unrelated"); n != 1 {
		t.Fatalf("getTransaction unrelated = %d, want 1: signatures below the saved cursor must not be processed again", n)
	}
}

func TestSolanaWatchDoesNotReuseTransaction(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	paid := createWatchTestOrder(t, srvURL, "trade-paid", model.OrderStatusSuccess, now.Add(-5*time.Minute), now.Add(5*time.Minute))
	model.Db.Model(&paid).Updates(map[string]any{"ref_hash": "pay-1"})
	waiting := createWatchTestOrder(t, srvURL, "trade-waiting", model.OrderStatusWaiting, now.Add(-2*time.Minute), now.Add(8*time.Minute))
	rpc.set(watchTestATA(t), watchTestPayment("pay-1", 5000, now.Add(-30*time.Second)))

	s.watchPayments(t.Context())

	if got := getTestOrder(t, waiting.TradeId); got.Status != model.OrderStatusWaiting {
		t.Fatalf("waiting order status = %d, want waiting: pay-1 already belongs to another order", got.Status)
	}
	if n := rpc.count("getTransaction pay-1"); n != 0 {
		t.Fatalf("getTransaction pay-1 = %d, want 0 for a used transaction", n)
	}
}

func TestSolanaMarkConfirmingOnceGuards(t *testing.T) {
	initTaskTestEnv(t, nil)
	now := time.Now()
	newOrder := func(tradeID string, status int) model.Order {
		return createWatchTestOrder(t, "http://127.0.0.1", tradeID, status, now.Add(-time.Minute), now.Add(9*time.Minute))
	}
	// 读出后被取消：不能复活
	canceled := newOrder("trade-canceled", model.OrderStatusWaiting)
	model.Db.Model(&model.Order{}).Where("id = ?", canceled.ID).Update("status", model.OrderStatusCanceled)
	if ok, err := canceled.MarkConfirmingOnce(1, "from", "hash-1", now, decimalOf(t, watchTestAmount)); ok || err != nil {
		t.Fatalf("MarkConfirmingOnce on a canceled order = %v %v, want false", ok, err)
	}
	if got := getTestOrder(t, canceled.TradeId); got.Status != model.OrderStatusCanceled {
		t.Fatalf("canceled order status = %d", got.Status)
	}

	// 同一交易只能让一个订单进入确认中
	first := newOrder("trade-first", model.OrderStatusWaiting)
	second := newOrder("trade-second", model.OrderStatusExpired)
	if ok, err := first.MarkConfirmingOnce(1, "from", "hash-2", now, decimalOf(t, watchTestAmount)); !ok || err != nil {
		t.Fatalf("first MarkConfirmingOnce = %v %v, want true", ok, err)
	}
	if ok, err := second.MarkConfirmingOnce(1, "from", "hash-2", now, decimalOf(t, watchTestAmount)); ok || err != nil {
		t.Fatalf("second MarkConfirmingOnce with the same hash = %v %v, want false", ok, err)
	}
	if got := getTestOrder(t, second.TradeId); got.Status != model.OrderStatusExpired || got.RefHash != "" {
		t.Fatalf("second order = status %d ref %q, want untouched", got.Status, got.RefHash)
	}
}

func TestSolanaWatchFailsOverWhenPollNodeRateLimited(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	rpc.set(watchTestATA(t), watchTestPayment("pay-1", 5000, now.Add(-30*time.Second)))
	rpc.setStatus("/poll", http.StatusTooManyRequests)
	createWatchTestIdleOrder(t, srvURL, now)

	s.watchPayments(t.Context())
	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming {
		t.Fatalf("order status = %d, want confirming through the main node", got.Status)
	}

	// 冷却期间直接用主节点
	s.watchPayments(t.Context())
	if n := rpc.count("/poll getSignaturesForAddress"); n != 1 {
		t.Fatalf("poll node getSignaturesForAddress = %d during cooldown, want 1", n)
	}
	if n := rpc.count("/main getSignaturesForAddress"); n != 2 {
		t.Fatalf("main node getSignaturesForAddress = %d, want 2", n)
	}
}

func TestSolanaWatchPagesThroughSignatures(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-2*time.Minute), now.Add(8*time.Minute))
	ata := watchTestATA(t)

	// 付款之后又有 150 笔失败交易，付款在第二页（首页 20 条，之后每页 1000 条）
	sigs := make([]watchTestSig, 0)
	for i := 150; i > 0; i-- {
		sigs = append(sigs, watchTestSig{sig: fmt.Sprintf("noise-%d", i), slot: int64(5000 + i), blockTime: now.Add(-10 * time.Second).Unix(), err: `{"InstructionError":[0,"InvalidAccountData"]}`})
	}
	sigs = append(sigs, watchTestPayment("pay-1", 5000, now.Add(-30*time.Second)))
	rpc.set(ata, sigs...)

	s.watchPayments(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" {
		t.Fatalf("order = status %d ref %q, want the payment on the second page", got.Status, got.RefHash)
	}
	if n := rpc.count("/poll getSignaturesForAddress"); n != 2 {
		t.Fatalf("getSignaturesForAddress pages = %d, want 2", n)
	}
}

func TestSolanaWatchRecentlyExpiredOrder(t *testing.T) {
	tests := []struct {
		name    string
		paidAgo time.Duration
		want    int
	}{
		{name: "paid before expiry", paidAgo: 90 * time.Second, want: model.OrderStatusConfirming},
		{name: "paid after expiry", paidAgo: 30 * time.Second, want: model.OrderStatusExpired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rpc := newWatchTestRPC(t)
			s, srvURL := initWatchTest(t, rpc)
			now := time.Now()
			order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusExpired, now.Add(-11*time.Minute), now.Add(-time.Minute))
			rpc.set(watchTestATA(t), watchTestPayment("pay-1", 5000, now.Add(-tt.paidAgo)))

			s.watchPayments(t.Context())

			if got := getTestOrder(t, order.TradeId); got.Status != tt.want {
				t.Fatalf("order status = %d, want %d", got.Status, tt.want)
			}
		})
	}
}

func TestSolanaWatchFindsNonAssociatedTokenAccount(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	const extra = "ExtraTokenAccount111111111111111111111111111"
	rpc.extra[watchTestOwner] = extra
	rpc.set(extra, watchTestPayment("pay-1", 5000, now.Add(-30*time.Second)))

	s.watchPayments(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" {
		t.Fatalf("order = status %d ref %q, want the payment to a non-associated token account", got.Status, got.RefHash)
	}
}

func TestSolanaScanModeSwitchesDetection(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	createWatchTestOrder(t, srvURL, "trade-watch-1", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))

	// address：整链扫描不取高度、不回溯
	s.syncSlotForward(t.Context())
	s.lookbackSlots(t.Context())
	if s.slotQueue.Len() != 0 || rpc.count("/main getSlot") != 0 {
		t.Fatalf("block scan ran in address mode: queue %d getSlot %d", s.slotQueue.Len(), rpc.count("/main getSlot"))
	}

	// block：不按账户检测
	model.Db.Model(&model.Conf{}).Where("k = ?", model.SolanaScanMode).Update("v", "block")
	model.RefreshC()
	s.watchPayments(t.Context())
	if n := rpc.count("/poll getSignaturesForAddress") + rpc.count("/main getTokenAccountsByOwner"); n != 0 {
		t.Fatalf("address detection made %d RPC calls in block mode", n)
	}
}
