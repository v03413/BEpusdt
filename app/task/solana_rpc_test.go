package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

// 付费节点的 API Key 在地址路径里，错误信息只能带主机名
func TestSolanaRPCErrorsHideEndpointPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	initTaskTestEnv(t, map[model.ConfKey]string{
		model.RpcEndpointSolana:     "http://127.0.0.1:1/v2/secret-main",
		model.RpcEndpointSolanaPoll: srv.URL + "/v2/secret-poll",
	})
	solanaCooldown.Clear()
	s := newSolana()

	_, err := s.rpcCall(t.Context(), "getSignaturesForAddress", "account")
	if err == nil {
		t.Fatal("expected an error from both nodes")
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("error leaks the endpoint path: %v", err)
	}
	if !strings.Contains(err.Error(), "127.0.0.1:1") || !strings.Contains(err.Error(), "rate limited") {
		t.Fatalf("error = %v, want both nodes named by host", err)
	}
}
