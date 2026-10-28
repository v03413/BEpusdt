package task

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/model"
)

// Solana RPC 分两个角色：主节点（rpc_endpoint_solana）负责取交易、确认状态、代币账户列表和区块；
// 轮询节点（rpc_endpoint_solana_poll，为空时同主节点）负责高频的 getSignaturesForAddress。
// 每个请求先发首选节点，失败再换另一个；被限流的节点按方法冷却一段时间，期间直接用另一个。
// 每次尝试单独限时，并给后面的节点留出时间：一个节点卡住（例如响应体迟迟不结束）不会耗尽整个请求的时限。
// 付费节点的 API Key 通常在地址路径里，日志和错误里只出现主机名。

const (
	solanaRateLimitCooldown = 5 * time.Second
	solanaAttemptTimeout    = 8 * time.Second // 单个节点单次请求的最长时间
)

// solanaCooldown 记录被限流的节点和方法，key: 节点地址 + " " + 方法，value: 冷却结束时间
var solanaCooldown sync.Map

var errSolanaNullResult = errors.New("result is null")

func solanaEndpoints(method string) []string {
	main := strings.TrimSpace(model.GetC(model.RpcEndpointSolana))
	poll := strings.TrimSpace(model.GetC(model.RpcEndpointSolanaPoll))
	if poll == "" || poll == main {

		return []string{main}
	}
	if method == "getSignaturesForAddress" {

		return []string{poll, main}
	}

	return []string{main, poll}
}

func (s *solana) rpcCall(ctx context.Context, method string, params ...any) (gjson.Result, error) {

	return s.rpcCallAccept(ctx, method, nil, params...)
}

// rpcCallAccept 依次尝试各节点；accept 不为空时，结果不被接受（例如节点还没有这笔交易，返回 null）也换下一个节点
func (s *solana) rpcCallAccept(ctx context.Context, method string, accept func(gjson.Result) error, params ...any) (gjson.Result, error) {
	post, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": 1, "method": method, "params": params})
	if err != nil {

		return gjson.Result{}, err
	}

	var errs []error
	endpoints := solanaEndpoints(method)
	for i, endpoint := range endpoints {
		if err := ctx.Err(); err != nil {

			return gjson.Result{}, errors.Join(append(errs, err)...)
		}
		if until, ok := solanaCooldown.Load(endpoint + " " + method); ok && time.Now().Before(until.(time.Time)) {
			errs = append(errs, fmt.Errorf("%s %s: rate limited, cooling down", method, solanaHost(endpoint)))

			continue
		}

		attemptCtx, cancel := solanaAttemptContext(ctx, len(endpoints)-i)
		result, err := s.rpcPost(attemptCtx, endpoint, method, post)
		cancel()
		if err == nil && accept != nil {
			if err = accept(result); err != nil {
				err = fmt.Errorf("%s %s: %w", method, solanaHost(endpoint), err)
			}
		}
		if err == nil {

			return result, nil
		}
		errs = append(errs, err)
	}

	return gjson.Result{}, errors.Join(errs...)
}

// solanaAttemptContext 返回单次尝试的时限：不超过 solanaAttemptTimeout，调用方有截止时间时按剩余节点数平分剩余时间
func solanaAttemptContext(ctx context.Context, left int) (context.Context, context.CancelFunc) {
	timeout := solanaAttemptTimeout
	if deadline, ok := ctx.Deadline(); ok && left > 1 {
		timeout = min(timeout, time.Until(deadline)/time.Duration(left))
	}

	return context.WithTimeout(ctx, timeout)
}

func (s *solana) rpcPost(ctx context.Context, endpoint, method string, post []byte) (gjson.Result, error) {
	host := solanaHost(endpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(post))
	if err != nil {

		return gjson.Result{}, fmt.Errorf("%s %s: invalid endpoint", method, host)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {

		return gjson.Result{}, fmt.Errorf("%s %s: %w", method, host, solanaRedactError(err))
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		solanaCooldown.Store(endpoint+" "+method, time.Now().Add(solanaRetryAfter(resp.Header.Get("Retry-After"))))

		return gjson.Result{}, fmt.Errorf("%s %s: rate limited (429)", method, host)
	}
	if resp.StatusCode != http.StatusOK {

		return gjson.Result{}, fmt.Errorf("%s %s: response status code %d", method, host, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return gjson.Result{}, fmt.Errorf("%s %s: %w", method, host, solanaRedactError(err))
	}
	if !gjson.ValidBytes(body) {

		return gjson.Result{}, fmt.Errorf("%s %s: invalid JSON response", method, host)
	}

	data := gjson.ParseBytes(body)
	if rpcErr := data.Get("error"); rpcErr.Exists() {

		return gjson.Result{}, fmt.Errorf("%s %s: RPC error: %s", method, host, rpcErr.Raw)
	}
	if !data.Get("result").Exists() {

		return gjson.Result{}, fmt.Errorf("%s %s: missing result", method, host)
	}

	return data.Get("result"), nil
}

func solanaRetryAfter(header string) time.Duration {
	if sec, err := strconv.Atoi(strings.TrimSpace(header)); err == nil && sec > 0 && sec <= 60 {

		return time.Duration(sec) * time.Second
	}

	return solanaRateLimitCooldown
}

// solanaHost 只保留主机名，避免 API Key 进入日志
func solanaHost(endpoint string) string {
	if u, err := url.Parse(endpoint); err == nil && u.Host != "" {

		return u.Host
	}

	return "solana-rpc"
}

// solanaRedactError 去掉 net/http 错误里附带的完整请求地址
func solanaRedactError(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {

		return urlErr.Err
	}

	return err
}

func solanaRequireNonNull(result gjson.Result) error {
	if result.Type == gjson.Null {

		return errSolanaNullResult
	}

	return nil
}
