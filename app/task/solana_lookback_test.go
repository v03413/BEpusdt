package task

import (
	"testing"
	"time"

	"github.com/v03413/bepusdt/app/model"
)

func createSolanaLookbackOrder(t *testing.T, tradeID string, status int) model.Order {
	t.Helper()

	return createTestOrder(t, model.Order{
		OrderId:   "merchant-" + tradeID,
		TradeId:   tradeID,
		TradeType: model.UsdcSolana,
		Rate:      "7.00",
		Amount:    "1.23",
		Money:     "8.61",
		Address:   "GpsvqTPP7triXkqH8mVtXEM6at9DJ41UHy5dYPTFe16H",
		Status:    status,
		ApiType:   model.OrderApiTypeEpusdt,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}, time.Now().Add(-time.Minute))
}

// 上一笔订单回溯到一半就取消了；一周后新订单到来时，不能接着扫上一笔订单留下的旧 slot
func TestSolanaLookbackDropsRangeOfFinishedOrders(t *testing.T) {
	initTaskTestEnv(t, nil)
	lookbackDone.Clear()

	canceled := createSolanaLookbackOrder(t, "lookback-canceled", model.OrderStatusCanceled)
	waiting := createSolanaLookbackOrder(t, "lookback-waiting", model.OrderStatusWaiting)
	lookbackDone.Store(waiting.ID, struct{}{}) // 新订单的回溯不在本测试范围，避免请求区块高度接口

	s := newSolana()
	s.lookbackNext, s.lookbackEnd, s.lookbackOrderIDs = 459552217, 459552221, []int64{canceled.ID}
	s.lookbackSlots(t.Context())

	select {
	case slot := <-s.slotQueue.Out:
		t.Fatalf("queued slot %d from the range of a canceled order", slot)
	case <-time.After(200 * time.Millisecond):
	}
	if s.lookbackNext != 0 || s.lookbackEnd != 0 || s.lookbackOrderIDs != nil {
		t.Fatalf("lookback progress kept: next = %d, end = %d, orders = %v", s.lookbackNext, s.lookbackEnd, s.lookbackOrderIDs)
	}
}

func TestSolanaLookbackResumesRangeOfPendingOrders(t *testing.T) {
	initTaskTestEnv(t, nil)
	lookbackDone.Clear()

	waiting := createSolanaLookbackOrder(t, "lookback-waiting", model.OrderStatusWaiting)

	const start, end = 462296900, 462296904
	s := newSolana()
	s.lookbackNext, s.lookbackEnd, s.lookbackOrderIDs = start, end, []int64{waiting.ID}
	s.lookbackSlots(t.Context())

	for want := start; want <= end; want++ {
		select {
		case got := <-s.slotQueue.Out:
			if got != want {
				t.Fatalf("queued slot = %d, want %d", got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("slot %d not queued", want)
		}
	}
	if s.lookbackNext != 0 {
		t.Fatalf("lookbackNext = %d after the range was queued, want 0", s.lookbackNext)
	}
	if _, done := lookbackDone.Load(waiting.ID); !done {
		t.Fatalf("order %d not marked as looked back", waiting.ID)
	}
}
