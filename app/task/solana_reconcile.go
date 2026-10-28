package task

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"

	"filippo.io/edwards25519"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
)

// 订单过期后，按收款账户的链上历史对账一次。
// 扫块每个 slot 只有一次机会，跳块、放弃重试的 slot、重启、回溯范围不完整都会让有效期内的付款漏掉；
// 收款账户的签名列表是独立的链上历史，可以把这些付款找回来。

const (
	solanaReconcileDelay     = time.Minute   // 过期后等待交易 finalized 再对账
	solanaReconcileWindow    = 2 * time.Hour // 只对最近过期的订单对账
	solanaSignaturePageLimit = 1000

	solanaAssociatedTokenProgram = "ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL"
)

// solanaReconciled 记录已完成对账的订单 ID；RPC 失败或代币账户列表不完整的订单不记录，下一轮重试
var solanaReconciled sync.Map // key: int64 order ID, value: struct{}

// solanaReconcileProgress 记录未完成对账的订单已经下载核对过的交易，重试时不重复下载。
// 每次重试都重新列出全部账户的签名：节点一时返回的空列表或缺项不能当作该账户已经查完
var solanaReconcileProgress sync.Map // key: int64 order ID, value: *solanaReconcileState

type solanaReconcileState struct {
	examined map[string]bool // 已下载核对、与该订单不匹配的交易签名
	warned   bool            // 已记过"代币账户列表不完整"
}

// errSolanaTokenAccountsIncomplete 表示没能列出钱包的全部代币账户，只查了已知的；没找到付款时不能算对账完成
var errSolanaTokenAccountsIncomplete = errors.New("token account list incomplete")

func (s *solana) reconcileExpiredOrders(ctx context.Context) {
	now := time.Now()
	var orders []model.Order
	model.Db.Where("status = ? and trade_type in (?)", model.OrderStatusExpired, model.GetNetworkTrades(conf.Solana)).
		Where("expired_at > ? and expired_at < ?", now.Add(-solanaReconcileWindow), now.Add(-solanaReconcileDelay)).
		Find(&orders)

	current := make(map[int64]bool, len(orders))
	for _, o := range orders {
		current[o.ID] = true
	}
	solanaReconcileProgress.Range(func(id, _ any) bool {
		if !current[id.(int64)] {
			solanaReconcileProgress.Delete(id)
		}

		return true
	})

	for _, o := range orders {
		if ctx.Err() != nil {

			return
		}
		if _, done := solanaReconciled.Load(o.ID); done {

			continue
		}

		v, _ := solanaReconcileProgress.LoadOrStore(o.ID, &solanaReconcileState{examined: make(map[string]bool)})
		state := v.(*solanaReconcileState)
		if err := s.reconcileOrder(ctx, o, state); err != nil {
			if !errors.Is(err, errSolanaTokenAccountsIncomplete) {
				log.Task.Warn(fmt.Sprintf("Solana 订单对账失败，稍后重试 %s：%v", o.TradeId, err))
			} else if !state.warned && model.GetC(model.RpcEndpointSolanaPoll) != "" { // 只配置一个节点时（如 publicnode 免费节点）该方法通常不可用，不记日志
				state.warned = true
				log.Task.Warn(fmt.Sprintf("Solana 订单 %s 已查过已知代币账户，未找到付款；钱包的代币账户列表暂时查不全，稍后重试：%v", o.TradeId, err))
			}

			continue
		}

		solanaReconciled.Store(o.ID, struct{}{})
		solanaReconcileProgress.Delete(o.ID)
	}
}

func (s *solana) reconcileOrder(ctx context.Context, o model.Order, state *solanaReconcileState) error {
	accounts, listErr, err := s.tokenAccounts(ctx, orderMatchAddress(o), model.GetTradeContract(o.TradeType))
	if err != nil {

		return err
	}

	for _, account := range accounts {
		signatures, err := s.orderSignatures(ctx, account, o)
		if err != nil {

			return err
		}

		for _, signature := range signatures {
			if state.examined[signature] {

				continue
			}

			var used int64
			if err := model.Db.Model(&model.Order{}).Where("ref_hash = ?", signature).Count(&used).Error; err != nil {

				return err
			}
			if used > 0 {

				continue
			}

			trans, err := s.rpcCallAccept(ctx, "getTransaction", solanaRequireTransaction, signature,
				map[string]any{"encoding": "json", "maxSupportedTransactionVersion": 1, "commitment": "finalized"})
			if err != nil {

				return err
			}

			timestamp := time.Unix(trans.Get("blockTime").Int(), 0)
			for _, t := range s.parseTransaction(trans, int(trans.Get("slot").Int()), timestamp) {
				if !orderTransferMatch(o, t) {

					continue
				}
				ok, err := solanaMarkConfirming(&o, t)
				if err != nil {

					return err
				}
				if !ok { // 订单状态已变，或交易已被其他订单使用；状态已变时后面的匹配也不会更新

					continue
				}

				log.Task.Info(fmt.Sprintf("Solana 对账找回付款 %s：%s", o.TradeId, t.TxHash))

				return nil
			}
			state.examined[signature] = true
		}
	}
	if listErr != nil {

		return fmt.Errorf("%w: %v", errSolanaTokenAccountsIncomplete, listErr)
	}

	return nil
}

// tokenAccounts 返回钱包在该 mint 下的代币账户：总是包含关联代币账户（ATA），节点支持 getTokenAccountsByOwner 时再加上其他代币账户。
// 部分节点（如 publicnode 免费节点）不提供 getTokenAccountsByOwner，此时只能先查 ATA（绝大多数付款都转入 ATA），listErr 不为空，
// 表示列表可能不完整
func (s *solana) tokenAccounts(ctx context.Context, owner, mint string) (accounts []string, listErr error, err error) {
	ata, err := solanaAssociatedTokenAccount(owner, mint)
	if err != nil {

		return nil, nil, err
	}

	accounts = []string{ata}
	extra, listErr := s.listTokenAccounts(ctx, owner, mint)
	for _, account := range extra {
		if account != ata {
			accounts = append(accounts, account)
		}
	}

	return accounts, listErr, nil
}

// listTokenAccounts 用 getTokenAccountsByOwner 列出钱包在该 mint 下的全部代币账户
func (s *solana) listTokenAccounts(ctx context.Context, owner, mint string) ([]string, error) {
	result, err := s.rpcCall(ctx, "getTokenAccountsByOwner", owner, map[string]any{"mint": mint},
		map[string]any{"encoding": "base64", "dataSlice": map[string]int{"offset": 0, "length": 0}, "commitment": "finalized"})
	if err != nil {

		return nil, err
	}

	accounts := make([]string, 0)
	for _, account := range result.Get("value.#.pubkey").Array() {
		accounts = append(accounts, account.String())
	}

	return accounts, nil
}

// solanaAssociatedTokenAccount 按 Associated Token Account 程序的规则推导地址：
// seeds = [钱包, Token 程序, mint]，从 bump 255 往下找第一个不在 ed25519 曲线上的哈希
func solanaAssociatedTokenAccount(owner, mint string) (string, error) {
	seeds := [][]byte{base58.Decode(owner), base58.Decode(conf.SolSplToken), base58.Decode(mint)}
	for _, seed := range seeds {
		if len(seed) != 32 {

			return "", fmt.Errorf("invalid Solana address: %s / %s", owner, mint)
		}
	}

	program := base58.Decode(solanaAssociatedTokenProgram)
	for bump := 255; bump >= 0; bump-- {
		h := sha256.New()
		for _, seed := range seeds {
			h.Write(seed)
		}
		h.Write([]byte{byte(bump)})
		h.Write(program)
		h.Write([]byte("ProgramDerivedAddress"))

		addr := h.Sum(nil)
		if _, err := new(edwards25519.Point).SetBytes(addr); err != nil {

			return base58.Encode(addr), nil
		}
	}

	return "", fmt.Errorf("no associated token account for %s / %s", owner, mint)
}

// orderSignatures 返回代币账户在订单有效期内执行成功的交易签名；签名按时间倒序分页返回
func (s *solana) orderSignatures(ctx context.Context, account string, o model.Order) ([]string, error) {
	var result []string
	opts := map[string]any{"limit": solanaSignaturePageLimit, "commitment": "finalized"}
	for {
		page, err := s.rpcCall(ctx, "getSignaturesForAddress", account, opts)
		if err != nil {

			return nil, err
		}

		items := page.Array()
		for _, itm := range items {
			if !itm.Get("blockTime").Exists() || itm.Get("blockTime").Type == gjson.Null {

				continue
			}

			at := time.Unix(itm.Get("blockTime").Int(), 0)
			if !o.CreatedAt.Before(at) { // 已早于订单创建时间

				return result, nil
			}
			if o.ExpiredAt.After(at) && solanaTxSucceeded(itm.Get("err")) {
				result = append(result, itm.Get("signature").String())
			}
		}
		if len(items) < solanaSignaturePageLimit {

			return result, nil
		}

		opts["before"] = items[len(items)-1].Get("signature").String()
	}
}
