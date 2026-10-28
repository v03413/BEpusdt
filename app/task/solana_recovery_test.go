package task

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
)

// 按账户检测在节点故障、并发状态变更时的恢复：慢请求、列表缺项、账户列表不完整、过期与匹配并发、确认状态为空

func TestSolanaExpiryDoesNotOverwriteMatchedPayment(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-expiry", model.OrderStatusWaiting, now.Add(-11*time.Minute), now.Add(-time.Second))
	// 过期任务已读出等待支付的订单，还没写回
	stale := getTestOrder(t, order.TradeId)
	canceling := getTestOrder(t, order.TradeId)

	rpc.set(watchTestATA(t), watchTestPayment("pay-1", 5000, now.Add(-2*time.Second)))
	watchRound(t, s)

	if ok, err := stale.SetExpired(); ok || err != nil {
		t.Fatalf("SetExpired on a matched order = %v %v, want false", ok, err)
	}
	if err := canceling.SetCanceled(); !errors.Is(err, model.ErrOrderStatusChanged) {
		t.Fatalf("SetCanceled on a matched order = %v, want ErrOrderStatusChanged", err)
	}
	got := getTestOrder(t, order.TradeId)
	if got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" || got.RefBlockNum != 5000 {
		t.Fatalf("order = status %d ref %q slot %d, want the matched payment kept", got.Status, got.RefHash, got.RefBlockNum)
	}

	// 过期任务只让仍在等待支付的订单过期
	waiting := createWatchTestOrder(t, srvURL, "trade-waiting", model.OrderStatusWaiting, now.Add(-11*time.Minute), now.Add(-time.Second))
	expireWaitingOrders()
	if got := getTestOrder(t, waiting.TradeId); got.Status != model.OrderStatusExpired {
		t.Fatalf("waiting order status = %d, want expired", got.Status)
	}
	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-1" {
		t.Fatalf("matched order = status %d ref %q after expireWaitingOrders", got.Status, got.RefHash)
	}
}

func TestSolanaCancelOnlyWaitingOrder(t *testing.T) {
	initTaskTestEnv(t, nil)
	now := time.Now()
	o := createWatchTestOrder(t, "http://127.0.0.1", "trade-cancel", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	if err := o.SetCanceled(); err != nil {
		t.Fatalf("SetCanceled = %v", err)
	}
	if got := getTestOrder(t, o.TradeId); got.Status != model.OrderStatusCanceled || o.Status != model.OrderStatusCanceled {
		t.Fatalf("order status = %d (in memory %d), want canceled", got.Status, o.Status)
	}
}

func TestSolanaWatchProcessesKnownSignatureMissingFromList(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-missing", model.OrderStatusWaiting, now.Add(-3*time.Minute), now.Add(7*time.Minute))
	ata := watchTestATA(t)
	stuck := watchTestSig{sig: "stuck-pay", slot: 4700, blockTime: now.Add(-90 * time.Second).Unix()} // 两个节点都还查不到
	noise := watchTestSig{sig: "new-noise", slot: 5100, blockTime: now.Add(-time.Second).Unix(), fixture: watchTestOther}
	rpc.set(ata, noise, stuck)

	watchRound(t, s)
	if c := getWatchTestCursor(t, ata); c.Slot != 4700 || !strings.Contains(c.Pending, stuck.sig) {
		t.Fatalf("cursor = %+v, want held at 4700 with the stuck signature pending", c)
	}

	// 另一个后端的列表暂时缺了这笔签名：低水位不能越过它
	rpc.set(ata, noise)
	s.watch.mu.Lock()
	s.watch.signatures[stuck.sig].retryAt = time.Time{}
	s.watch.mu.Unlock()
	watchRound(t, s)
	if c := getWatchTestCursor(t, ata); c.Slot != 4700 {
		t.Fatalf("cursor slot = %d while the signature is missing from the list, want 4700", c.Slot)
	}
	if n := rpc.count("getTransaction " + stuck.sig); n != 4 {
		t.Fatalf("getTransaction %s = %d, want it fetched by signature while missing from the list (2 per round)", stuck.sig, n)
	}

	// 节点有了这笔交易，但列表仍然缺它：按签名直接下载并入账
	stuck.fixture = watchTestPay
	rpc.mu.Lock()
	rpc.txs[stuck.sig] = stuck
	rpc.mu.Unlock()
	s.watch.mu.Lock()
	s.watch.signatures[stuck.sig].retryAt = time.Time{}
	s.watch.mu.Unlock()
	watchRound(t, s)

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != stuck.sig {
		t.Fatalf("order = status %d ref %q, want the payment missing from the list", got.Status, got.RefHash)
	}
	if c := getWatchTestCursor(t, ata); c.Slot != 5100-solanaWatchSlotOverlap || c.Pending != "" {
		t.Fatalf("cursor = %+v, want slot %d with nothing pending", c, 5100-solanaWatchSlotOverlap)
	}
}

func TestSolanaWatchRestoresPendingSignaturesAfterRestart(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-restart", model.OrderStatusWaiting, now.Add(-3*time.Minute), now.Add(7*time.Minute))
	ata := watchTestATA(t)
	stuck := watchTestSig{sig: "stuck-pay", slot: 4700, blockTime: now.Add(-90 * time.Second).Unix()}
	noise := watchTestSig{sig: "new-noise", slot: 5100, blockTime: now.Add(-time.Second).Unix(), fixture: watchTestOther}
	rpc.set(ata, noise, stuck)
	watchRound(t, s)

	// 重启：内存状态全部丢失；重启后的列表缺了这笔签名，节点已能查到交易
	rpc.set(ata, noise)
	stuck.fixture = watchTestPay
	rpc.mu.Lock()
	rpc.txs[stuck.sig] = stuck
	rpc.mu.Unlock()
	restarted := newSolana()
	watchRound(t, &restarted)

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != stuck.sig {
		t.Fatalf("order = status %d ref %q after restart, want the pending payment processed", got.Status, got.RefHash)
	}
	if c := getWatchTestCursor(t, ata); c.Pending != "" || c.Slot != 5100-solanaWatchSlotOverlap {
		t.Fatalf("cursor = %+v after restart, want slot %d with nothing pending", c, 5100-solanaWatchSlotOverlap)
	}
}

func TestSolanaWatchSlowTokenAccountListDoesNotBlockPolling(t *testing.T) {
	var polls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch gjson.GetBytes(body, "method").String() {
		case "getTokenAccountsByOwner": // 响应头已返回，响应体一直不结束
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		case "getSignaturesForAddress":
			polls.Add(1)
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":[]}`))
		default:
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":null}`))
		}
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.SolanaScanMode: solanaScanModeAddress, model.RpcEndpointSolana: srv.URL})
	solanaCooldown.Clear()
	s := newSolana()
	now := time.Now()
	createWatchTestOrder(t, srv.URL, "trade-slow-list", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))

	ctx, cancel := context.WithCancel(t.Context())
	defer func() {
		cancel()
		s.watch.inflight.Wait()
	}()

	start := time.Now()
	s.watchPayments(ctx)
	if d := time.Since(start); d > time.Second {
		t.Fatalf("watchPayments took %s, want it not to wait for the token account list", d)
	}
	waitFor(t, time.Second, func() bool { return polls.Load() > 0 }, "the known ATA was not polled while the token account list was slow")
}

func TestSolanaWatchSlowAccountDoesNotDelayOthers(t *testing.T) {
	usdc, _ := solanaAssociatedTokenAccount(watchTestOwner, model.GetTradeContract(model.UsdcSolana))
	usdt, _ := solanaAssociatedTokenAccount(watchTestOwner, model.GetTradeContract(model.UsdtSolana))
	var mu sync.Mutex
	polls := make(map[string]int)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch gjson.GetBytes(body, "method").String() {
		case "getSignaturesForAddress":
			account := gjson.GetBytes(body, "params.0").String()
			mu.Lock()
			polls[account]++
			mu.Unlock()
			if account == usdc { // 这个账户的查询卡住
				w.WriteHeader(http.StatusOK)
				w.(http.Flusher).Flush()
				<-r.Context().Done()

				return
			}
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":[]}`))
		default:
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"not supported"}}`))
		}
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.SolanaScanMode: solanaScanModeAddress, model.RpcEndpointSolana: srv.URL})
	solanaCooldown.Clear()
	s := newSolana()
	now := time.Now()
	createWatchTestOrder(t, srv.URL, "trade-usdc", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	usdtOrder := createWatchTestOrder(t, srv.URL, "trade-usdt", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	model.Db.Model(&usdtOrder).Update("trade_type", model.UsdtSolana)

	ctx, cancel := context.WithCancel(t.Context())
	defer func() {
		cancel()
		s.watch.inflight.Wait()
	}()

	count := func(account string) int {
		mu.Lock()
		defer mu.Unlock()

		return polls[account]
	}
	running := func(account string) bool {
		st := s.watch.account(account)
		s.watch.mu.Lock()
		defer s.watch.mu.Unlock()

		return st.running
	}
	for round := 1; round <= 3; round++ {
		waitFor(t, time.Second, func() bool { return !running(usdt) }, "the healthy account's poll did not finish")
		s.watchPayments(ctx)
		waitFor(t, time.Second, func() bool { return count(usdt) >= round }, "the healthy account missed a round")
	}
	if n := count(usdc); n != 1 {
		t.Fatalf("slow account polled %d times while its first poll was still running, want 1", n)
	}
}

func TestSolanaRPCFallbackGetsTimeAfterSlowPrimary(t *testing.T) {
	var backup atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/main" {
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			<-r.Context().Done()

			return
		}
		backup.Add(1)
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"blockTime":1,"meta":{"err":null}}}`))
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL + "/main", model.RpcEndpointSolanaPoll: srv.URL + "/poll"})
	solanaCooldown.Clear()
	s := newSolana()

	ctx, cancel := context.WithTimeout(t.Context(), 400*time.Millisecond)
	defer cancel()
	result, err := s.rpcCallAccept(ctx, "getTransaction", solanaRequireTransaction, "pay-1", map[string]any{"commitment": "confirmed"})
	if err != nil || backup.Load() != 1 || result.Get("blockTime").Int() != 1 {
		t.Fatalf("rpcCallAccept = %s, %v; backup calls %d, want the backup answer", result.Raw, err, backup.Load())
	}
}

func TestSolanaConfirmUsesBackupWhenStatusNull(t *testing.T) {
	var backup atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/notify":
			_, _ = w.Write([]byte("ok"))
		case "/main":
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"context":{"slot":9000},"value":[null]}}`))
		default:
			backup.Add(1)
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"context":{"slot":9000},"value":[{"slot":8000,"err":null,"confirmationStatus":"finalized"}]}}`))
		}
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL + "/main", model.RpcEndpointSolanaPoll: srv.URL + "/poll"})
	solanaCooldown.Clear()
	s := newSolana()
	now := time.Now()
	order := createWatchTestOrder(t, srv.URL, "trade-confirm", model.OrderStatusConfirming, now.Add(-time.Minute), now.Add(9*time.Minute))
	model.Db.Model(&order).Updates(map[string]any{"ref_hash": "pay-1", "ref_block_num": 8000})

	s.tradeConfirmHandle(t.Context())

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusSuccess || backup.Load() != 1 {
		t.Fatalf("order status = %d, backup calls %d; want success through the backup node", got.Status, backup.Load())
	}
	waitOrderNotified(t, order.TradeId)
}

func TestSolanaReconcileRetriesWhenTokenAccountListFails(t *testing.T) {
	rpc := newWatchTestRPC(t)
	s, srvURL := initWatchTest(t, rpc)
	solanaReconciled.Clear()
	solanaReconcileProgress.Clear()
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-non-ata", model.OrderStatusExpired, now.Add(-15*time.Minute), now.Add(-5*time.Minute))
	const extra = "ExtraTokenAccount111111111111111111111111111"
	rpc.extra[watchTestOwner] = extra
	rpc.set(watchTestATA(t), watchTestSig{sig: "ata-noise", slot: 4000, blockTime: now.Add(-8 * time.Minute).Unix(), fixture: watchTestOther})
	rpc.set(extra, watchTestPayment("pay-non-ata", 5000, now.Add(-6*time.Minute)))

	// 主节点故障，轮询节点不提供 getTokenAccountsByOwner：只能查 ATA，不能算对账完成
	rpc.setStatus("/main", http.StatusServiceUnavailable)
	s.reconcileExpiredOrders(t.Context())
	if _, done := solanaReconciled.Load(order.ID); done {
		t.Fatal("order marked as reconciled with an incomplete token account list")
	}
	if n := rpc.count("getTransaction ata-noise"); n != 2 { // 主节点 503，轮询节点返回
		t.Fatalf("getTransaction ata-noise = %d, want the ATA checked once", n)
	}

	rpc.setStatus("/main", 0)
	s.reconcileExpiredOrders(t.Context())
	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusConfirming || got.RefHash != "pay-non-ata" {
		t.Fatalf("order = status %d ref %q, want the non-ATA payment found after recovery", got.Status, got.RefHash)
	}
	if n := rpc.count("getTransaction ata-noise"); n != 2 {
		t.Fatalf("getTransaction ata-noise = %d, want the checked ATA not downloaded again", n)
	}
	if _, done := solanaReconciled.Load(order.ID); !done {
		t.Fatal("order not marked as reconciled after the payment was found")
	}
}

func TestSolanaBlockScanDoesNotMatchOrdersInAddressMode(t *testing.T) {
	rpc := newWatchTestRPC(t)
	_, srvURL := initWatchTest(t, rpc)
	now := time.Now()
	order := createWatchTestOrder(t, srvURL, "trade-block", model.OrderStatusWaiting, now.Add(-time.Minute), now.Add(9*time.Minute))
	for notOrderQueue.Len() > 0 {
		<-notOrderQueue.Out
	}

	paid := transfer{Network: conf.Solana, TxHash: "pay-1", Amount: decimalOf(t, watchTestAmount), FromAddress: "payer", RecvAddress: watchTestOwner,
		Timestamp: now.Add(-10 * time.Second), TradeType: model.UsdcSolana, BlockNum: 5000}
	unrelated := paid
	unrelated.TxHash = "other-1"
	unrelated.Amount = decimalOf(t, "1.5")
	handleOrderTransfers([]transfer{paid, unrelated})

	if got := getTestOrder(t, order.TradeId); got.Status != model.OrderStatusWaiting || got.RefHash != "" {
		t.Fatalf("order = status %d ref %q, want untouched by the block scan in address mode", got.Status, got.RefHash)
	}
	select {
	case other := <-notOrderQueue.Out:
		if len(other) != 1 || other[0].TxHash != "other-1" {
			t.Fatalf("non-order transfers = %+v, want only other-1", other)
		}
	case <-time.After(time.Second):
		t.Fatal("unrelated transfer not sent to non-order notification")
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool, msg string) {
	t.Helper()

	for deadline := time.Now().Add(timeout); time.Now().Before(deadline); time.Sleep(10 * time.Millisecond) {
		if cond() {
			return
		}
	}
	t.Fatal(msg)
}
