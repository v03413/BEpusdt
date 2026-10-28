package task

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

const solanaBlockCleanedUpResponse = `{"jsonrpc":"2.0","error":{"code":-32001,"message":"Block 459552217 cleaned up, does not exist on node. First available block: 461715420"},"id":1}`

func TestValidateSolanaBlockResponseCleanedUp(t *testing.T) {
	if _, err := validateSolanaBlockResponse([]byte(solanaBlockCleanedUpResponse)); !errors.Is(err, errSolanaBlockCleanedUp) {
		t.Fatalf("error = %v, want errSolanaBlockCleanedUp", err)
	}
}

// 并发处理同一个 slot 时，较晚返回的普通错误不能在 BlockCleanedUp 已经终止后重新建立重试
func TestSolanaSlotParseCleanedUpWinsConcurrentFailure(t *testing.T) {
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	var calls atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(firstStarted)
			<-releaseFirst
			w.WriteHeader(http.StatusTooManyRequests)

			return
		}
		_, _ = w.Write([]byte(solanaBlockCleanedUpResponse))
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL})

	const slot = 459552217
	s := newSolana()
	t.Cleanup(s.stopSlotRetryTimers)
	done := make(chan struct{})
	go func() {
		s.slotParse(slot)
		close(done)
	}()
	<-firstStarted
	s.slotParse(slot)
	close(releaseFirst)
	<-done

	if got := s.pendingSlotRetryCount(); got != 0 {
		t.Fatalf("pending retries after concurrent cleaned up = %d, want 0", got)
	}
	s.retryMu.Lock()
	_, scheduled := s.retryScheduled[slot]
	s.retryMu.Unlock()
	if scheduled {
		t.Fatalf("concurrent failure rescheduled cleaned up slot")
	}
}

// 节点已清理的 slot 只请求一次就放弃，不留下待重试记录，向前扫描不会因此暂停；之前因为限流留下的重试也一并清除
func TestSolanaSlotParseGivesUpCleanedUpSlot(t *testing.T) {
	var cleaned atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cleaned.Load() {
			w.WriteHeader(http.StatusTooManyRequests)

			return
		}
		_, _ = w.Write([]byte(solanaBlockCleanedUpResponse))
	}))
	defer srv.Close()
	initTaskTestEnv(t, map[model.ConfKey]string{model.RpcEndpointSolana: srv.URL})

	const slot = 459552217
	s := newSolana()
	t.Cleanup(s.stopSlotRetryTimers)

	s.slotParse(slot)
	if got := s.pendingSlotRetryCount(); got != 1 {
		t.Fatalf("pending retries after 429 = %d, want 1", got)
	}
	s.expireSlotRetryTimer(slot)

	cleaned.Store(true)
	s.slotParse(slot)
	if got := s.pendingSlotRetryCount(); got != 0 {
		t.Fatalf("pending retries after cleaned up = %d, want 0: forward sync stays paused", got)
	}
	s.retryMu.Lock()
	_, scheduled := s.retryScheduled[slot]
	s.retryMu.Unlock()
	if scheduled {
		t.Fatalf("cleaned up slot still scheduled for retry")
	}
}
