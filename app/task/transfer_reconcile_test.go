package task

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
)

func TestGetReceivableOrdersIncludesRecentExpiredAndFailedOrders(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "receivable-orders.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	now := time.Now()
	createOrder := func(tradeID string, status int, expiredAt time.Time) {
		createdAt := model.Datetime(expiredAt.Add(-30 * time.Minute))
		order := model.Order{
			OrderId:      tradeID,
			TradeId:      tradeID,
			TradeType:    model.UsdtBep20,
			Crypto:       model.USDT,
			Amount:       "13.19",
			Money:        "100",
			Address:      "0x0000000000000000000000000000000000000019",
			MatchAddress: "0x0000000000000000000000000000000000000019",
			Status:       status,
			NotifyUrl:    "https://example.com/webhook",
			ConfirmedAt:  &now,
			ExpiredAt:    expiredAt,
			AutoTimeAt: model.AutoTimeAt{
				CreatedAt: &createdAt,
				UpdatedAt: &createdAt,
			},
		}
		if err := model.Db.Create(&order).Error; err != nil {
			t.Fatalf("create order %s: %v", tradeID, err)
		}
	}

	createOrder("recent-expired", model.OrderStatusExpired, now.Add(-time.Hour))
	createOrder("recent-failed", model.OrderStatusFailed, now.Add(-2*time.Hour))
	createOrder("old-expired", model.OrderStatusExpired, now.Add(-25*time.Hour))
	createOrder("success", model.OrderStatusSuccess, now.Add(-time.Hour))

	orders := getReceivableOrders()
	key := "0x0000000000000000000000000000000000000019" + string(model.UsdtBep20)
	got := orders[key]
	if len(got) != 2 {
		t.Fatalf("expected 2 recent recoverable orders, got %d", len(got))
	}

	ids := map[string]bool{}
	for _, order := range got {
		ids[order.TradeId] = true
	}
	if !ids["recent-expired"] || !ids["recent-failed"] {
		t.Fatalf("expected recent expired and failed orders, got %+v", ids)
	}
	if ids["old-expired"] || ids["success"] {
		t.Fatalf("unexpected non-recoverable orders included: %+v", ids)
	}
}

func TestGetOldestRecoverableOrderIncludesExpiredOrders(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "oldest-recoverable-order.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	now := time.Now()
	createdAt := model.Datetime(now.Add(-2 * time.Hour))
	confirmedAt := time.Unix(0, 0)
	order := model.Order{
		OrderId:      "expired-bsc",
		TradeId:      "expired-bsc",
		TradeType:    model.UsdtBep20,
		Crypto:       model.USDT,
		Amount:       "13.19",
		Money:        "100",
		Address:      "0x0000000000000000000000000000000000000019",
		MatchAddress: "0x0000000000000000000000000000000000000019",
		Status:       model.OrderStatusExpired,
		NotifyUrl:    "https://example.com/webhook",
		ConfirmedAt:  &confirmedAt,
		ExpiredAt:    now.Add(-time.Hour),
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create expired order: %v", err)
	}

	got, ok := getOldestRecoverableOrder(conf.Bsc)
	if !ok {
		t.Fatal("expected expired order to be recoverable")
	}
	if got.TradeId != "expired-bsc" {
		t.Fatalf("expected expired-bsc, got %s", got.TradeId)
	}
}

func TestEvmReconcileRecoverableOrdersQueuesHistoricalBlocks(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "evm-reconcile-queue.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	createdAt := model.Datetime(time.Now().Add(-2 * time.Minute))
	confirmedAt := time.Unix(0, 0)
	order := model.Order{
		OrderId:      "expired-bsc",
		TradeId:      "expired-bsc",
		TradeType:    model.UsdtBep20,
		Crypto:       model.USDT,
		Amount:       "13.19",
		Money:        "100",
		Address:      "0x0000000000000000000000000000000000000019",
		MatchAddress: "0x0000000000000000000000000000000000000019",
		Status:       model.OrderStatusExpired,
		NotifyUrl:    "https://example.com/webhook",
		ConfirmedAt:  &confirmedAt,
		ExpiredAt:    time.Now().Add(-time.Minute),
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create expired order: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode rpc request: %v", err)
		}
		if req["method"] != "eth_blockNumber" {
			t.Fatalf("unexpected rpc method: %v", req["method"])
		}

		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x64"}`))
	}))
	defer server.Close()

	model.SetK(model.RpcEndpointBsc, server.URL)
	scanner := evm{
		Network:        conf.Bsc,
		Client:         server.Client(),
		AvgBlockTime:   1,
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](context.Background(), 30),
	}

	scanner.reconcileRecoverableOrders(context.Background())

	select {
	case got := <-scanner.blockScanQueue.Out:
		if got.From <= 0 || got.To < got.From {
			t.Fatalf("unexpected queued range: %+v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("expected historical block range to be queued")
	}
}

func TestSyncBreakKeepsScanningForExpiredRecoverableOrders(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "recoverable-sync-break.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	createdAt := model.Datetime(time.Now().Add(-2 * time.Minute))
	confirmedAt := time.Unix(0, 0)
	order := model.Order{
		OrderId:      "recoverable-bsc",
		TradeId:      "recoverable-bsc",
		TradeType:    model.UsdtBep20,
		Crypto:       model.USDT,
		Amount:       "13.19",
		Money:        "100",
		Address:      "0x0000000000000000000000000000000000000019",
		MatchAddress: "0x0000000000000000000000000000000000000019",
		Status:       model.OrderStatusExpired,
		NotifyUrl:    "https://example.com/webhook",
		ConfirmedAt:  &confirmedAt,
		ExpiredAt:    time.Now().Add(-time.Minute),
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create expired order: %v", err)
	}

	if syncBreak(conf.Bsc, 0) {
		t.Fatal("expected EVM scanner to keep running for recoverable expired order")
	}
}

func TestTronSyncBreakKeepsScanningForExpiredRecoverableOrders(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "recoverable-tron-sync-break.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	createdAt := model.Datetime(time.Now().Add(-2 * time.Minute))
	confirmedAt := time.Unix(0, 0)
	order := model.Order{
		OrderId:      "recoverable-tron",
		TradeId:      "recoverable-tron",
		TradeType:    model.UsdtTrc20,
		Crypto:       model.USDT,
		Amount:       "13.19",
		Money:        "100",
		Address:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		MatchAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		Status:       model.OrderStatusExpired,
		NotifyUrl:    "https://example.com/webhook",
		ConfirmedAt:  &confirmedAt,
		ExpiredAt:    time.Now().Add(-time.Minute),
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create expired order: %v", err)
	}

	scanner := newTron()
	if scanner.syncBreak() {
		t.Fatal("expected Tron scanner to keep running for recoverable expired order")
	}
}
