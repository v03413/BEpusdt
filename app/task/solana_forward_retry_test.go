package task

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/v03413/bepusdt/app/model"
)

// 模拟 n 个 slot 正在等待重试
func (s *solana) setPendingSlotRetries(first, n int) {
	s.retryMu.Lock()
	defer s.retryMu.Unlock()

	for slot := first; slot < first+n; slot++ {
		s.retryAttempts[slot] = 1
	}
}

// 一个 slot 在等待重试时，向前扫描照常推进，之后的付款不用等它重试完
func TestSolanaForwardContinuesWhileFewSlotsRetry(t *testing.T) {
	var head atomic.Int64
	srv := newSolanaSlotTestServer(t, &head)
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL, model.BlockHeightMaxDiff: fmt.Sprint(solanaTestMaxDiff)})
	createSolanaWaitingOrder(t)

	const last = 458000000
	s := newSolana()
	s.lastSlotNum = last
	s.setPendingSlotRetries(last-3, 1)
	head.Store(last + 5)
	s.syncSlotForward(t.Context())

	if s.lastSlotNum != last+5 {
		t.Fatalf("lastSlotNum = %d, want %d: forward sync paused by one slot retry", s.lastSlotNum, last+5)
	}
	select {
	case got := <-s.slotQueue.Out:
		if got != last+1 {
			t.Fatalf("first queued slot = %d, want %d", got, last+1)
		}
	case <-time.After(time.Second):
		t.Fatalf("no slot queued while one slot retries")
	}
}

// 很多 slot 同时等待重试，说明节点整体故障，暂停向前扫描
func TestSolanaForwardPausesWhileManySlotsRetry(t *testing.T) {
	var head atomic.Int64
	srv := newSolanaSlotTestServer(t, &head)
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL, model.BlockHeightMaxDiff: fmt.Sprint(solanaTestMaxDiff)})
	createSolanaWaitingOrder(t)

	const last = 458000000
	s := newSolana()
	s.lastSlotNum = last
	s.setPendingSlotRetries(last-solanaForwardPauseRetries, solanaForwardPauseRetries)
	head.Store(last + 5)
	s.syncSlotForward(t.Context())

	if s.lastSlotNum != last || s.slotQueue.Len() != 0 {
		t.Fatalf("lastSlotNum = %d, queued = %d; want forward sync paused at %d", s.lastSlotNum, s.slotQueue.Len(), last)
	}
}
