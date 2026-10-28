package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
)

// 按收款代币账户检测付款（solana_scan_mode=address）。
// 逐块下载要处理链上所有区块，节点稍有限流就排队、重试、补扫；而我们只关心少数几个收款代币账户有没有新入账。
// 这里只查询有待支付订单的收款代币账户的新签名，只下载这些交易：
//   - 同一个代币账户供多个订单共用，每轮只查一次；
//   - 进度按账户记在 bep_solana_cursor（低水位 slot：小于它的签名都已处理），重启后从这里继续；
//   - 单笔交易取不到时单独退避重试，低水位停在它前面，更新的付款照常处理；
//   - 已发现但没处理完的签名记在账户的待处理表里（随低水位一起保存），低水位不越过它们；
//     签名列表暂时缺了它（负载均衡后面的节点索引不一致）也按签名单独下载，不依赖它重新出现在列表里；
//   - 识别到的转账在这里同步匹配订单，订单进入确认中之后低水位才前进，进程中途退出不会丢付款；
//   - 各账户独立检测，一个账户或代币账户列表的查询变慢，不耽误其他账户的下一轮。
// 订单只由这里和过期后的对账匹配；整链扫描此时只为 MQTT 和"其它通知"运行，不改订单。
// 确认（finalized）、回调和订单过期后的对账沿用原流程。

const (
	solanaScanModeAddress = "address"

	solanaWatchInterval        = 2 * time.Second
	solanaWatchExpiredGrace    = 2 * time.Minute  // 过期后继续检测，覆盖最后几秒付款、过期后才确认的交易
	solanaWatchCreatedMargin   = time.Minute      // 往订单创建时间之前多查一段，覆盖 blockTime 与本机时钟的误差
	solanaWatchSlotOverlap     = 150              // 低水位比最新签名落后约 1 分钟，覆盖负载均衡节点之间的索引差异
	solanaWatchFirstPage       = 20               // 首页条数；冷清的账户一页就够，不够再按 1000 条翻页
	solanaWatchMaxPages        = 20               // 每轮最多翻页数，超过说明账户异常繁忙
	solanaWatchFetchPerRound   = 10               // 每个账户每轮最多下载的交易数，从最新的开始
	solanaWatchPollTimeout     = 15 * time.Second // 单个账户一轮的总时长
	solanaWatchFetchTimeout    = 5 * time.Second  // 单笔交易一次下载（含主备节点）的最长时间，给同一账户后面的签名留出时间
	solanaWatchGiveUp          = 30 * time.Minute // 单笔交易持续取不到时放弃，订单过期后由对账再查
	solanaWatchSignatureTTL    = time.Hour        // 已处理签名在内存里保留的时间
	solanaWatchWarnAfter       = 30 * time.Second // 账户持续查询失败多久记 WARN
	solanaWatchErrorAfter      = 2 * time.Minute  // 账户持续查询失败多久记 ERROR
	solanaTokenAccountsTTL     = 30 * time.Second // 钱包名下代币账户列表的刷新间隔
	solanaTokenAccountsTimeout = 10 * time.Second // 刷新代币账户列表的最长时间
	solanaConfirmingAlertAfter = 5 * time.Minute  // 确认中超过这个时间仍未成功记 ERROR
)

type solanaWatch struct {
	mu            sync.Mutex
	accounts      map[string]*solanaAccountState   // key: 代币账户
	signatures    map[string]*solanaSignatureState // key: 签名
	tokenAccounts map[string]*solanaTokenAccounts  // key: 钱包 + " " + mint
	inflight      sync.WaitGroup                   // 进行中的账户检测和代币账户列表刷新

	confirmAlerted sync.Map // 已告警的确认中订单 ID
}

// solanaAccountState 只由该账户正在运行的那次检测读写（running 保证同一账户同时只有一次），running 本身受 solanaWatch.mu 保护
type solanaAccountState struct {
	running      bool
	cursor       model.SolanaCursor
	pending      map[string]int64 // 已发现、尚未处理完的签名 -> slot
	loaded       bool
	failures     int
	retryAt      time.Time
	failingSince time.Time
	alerted      int // 本次故障已记的级别：1 WARN，2 ERROR
}

type solanaSignatureState struct {
	done      bool
	attempts  int
	retryAt   time.Time
	firstSeen time.Time
}

type solanaTokenAccounts struct {
	accounts   []string
	refreshAt  time.Time
	refreshing bool
	listErr    bool
}

type solanaWatchTarget struct {
	account   string
	owner     string
	tradeType model.TradeType
	since     time.Time // 早于这个时间的签名与当前订单无关
}

func newSolanaWatch() *solanaWatch {
	return &solanaWatch{
		accounts:      make(map[string]*solanaAccountState),
		signatures:    make(map[string]*solanaSignatureState),
		tokenAccounts: make(map[string]*solanaTokenAccounts),
	}
}

func solanaAddressMode() bool {

	return model.GetC(model.SolanaScanMode) == solanaScanModeAddress
}

// solanaMonitorEnabled 判断是否开启了需要整链数据的监控功能：MQTT 推送或"其它通知"钱包
func solanaMonitorEnabled() bool {
	if mqttSubscribed(conf.Solana) {

		return true
	}

	var count int64
	model.Db.Model(&model.Wallet{}).
		Where("other_notify = ? and trade_type in (?)", model.WaOtherEnable, model.GetNetworkTrades(conf.Solana)).
		Count(&count)

	return count > 0
}

// watchPayments 为每个要检测的代币账户启动一次检测，不等它们结束：还在运行的账户本轮跳过，其他账户照常按间隔检测
func (s *solana) watchPayments(ctx context.Context) {
	if !solanaAddressMode() {

		return
	}

	for _, g := range s.watchGroups() {
		for _, account := range s.watchTokenAccounts(ctx, g) {
			s.startPoll(ctx, g.target(account))
		}
	}

	s.watch.prune(time.Now())
}

type solanaWatchGroup struct {
	owner     string
	tradeType model.TradeType
	since     time.Time // 早于这个时间的签名与当前订单无关
}

func (g solanaWatchGroup) target(account string) solanaWatchTarget {

	return solanaWatchTarget{account: account, owner: g.owner, tradeType: g.tradeType, since: g.since}
}

// watchGroups 按收款地址和币种汇总需要检测的订单：等待支付的订单，以及过期不到 solanaWatchExpiredGrace 的订单
func (s *solana) watchGroups() []solanaWatchGroup {
	trades := model.GetNetworkTrades(conf.Solana)
	if len(trades) == 0 {

		return nil
	}

	var orders []model.Order
	err := model.Db.Where("trade_type in (?)", trades).
		Where("(status = ? or (status = ? and expired_at > ?))", model.OrderStatusWaiting, model.OrderStatusExpired, time.Now().Add(-solanaWatchExpiredGrace)).
		Order("created_at asc").
		Find(&orders).Error
	if err != nil {
		log.Task.Warn("Solana 读取待检测订单失败：", err)

		return nil
	}

	groups := make([]solanaWatchGroup, 0)
	index := make(map[string]int)
	for _, o := range orders {
		if o.CreatedAt == nil {

			continue
		}

		key := orderMatchAddress(o) + " " + string(o.TradeType)
		since := o.CreatedAt.Time().Add(-solanaWatchCreatedMargin)
		if i, ok := index[key]; ok {
			if since.Before(groups[i].since) {
				groups[i].since = since
			}

			continue
		}
		index[key] = len(groups)
		groups = append(groups, solanaWatchGroup{owner: orderMatchAddress(o), tradeType: o.TradeType, since: since})
	}

	return groups
}

// watchTokenAccounts 返回钱包在该币种下已知的代币账户：推导出的 ATA（账户还不存在也查，第一笔付款通常会创建它），
// 加上 getTokenAccountsByOwner 列出过的其他代币账户。列表到了刷新时间就在后台刷新，不等它返回；新发现的账户立即开始检测
func (s *solana) watchTokenAccounts(ctx context.Context, g solanaWatchGroup) []string {
	mint := model.GetTradeContract(g.tradeType)
	key := g.owner + " " + mint

	s.watch.mu.Lock()
	entry, ok := s.watch.tokenAccounts[key]
	if !ok {
		ata, err := solanaAssociatedTokenAccount(g.owner, mint)
		if err != nil {
			s.watch.mu.Unlock()
			log.Task.Warn("Solana 收款地址无效，跳过检测：", err)

			return nil
		}
		entry = &solanaTokenAccounts{accounts: []string{ata}}
		s.watch.tokenAccounts[key] = entry
	}
	accounts := append([]string(nil), entry.accounts...)
	refresh := !entry.refreshing && !time.Now().Before(entry.refreshAt)
	if refresh {
		entry.refreshing = true
		s.watch.inflight.Add(1)
	}
	s.watch.mu.Unlock()

	if refresh {
		go func() {
			defer s.watch.inflight.Done()

			s.refreshTokenAccounts(ctx, g, key, mint)
		}()
	}

	return accounts
}

// refreshTokenAccounts 用 getTokenAccountsByOwner 补充钱包的其他代币账户；节点不支持或查询失败时沿用已知账户
func (s *solana) refreshTokenAccounts(ctx context.Context, g solanaWatchGroup, key, mint string) {
	listCtx, cancel := context.WithTimeout(ctx, solanaTokenAccountsTimeout)
	extra, listErr := s.listTokenAccounts(listCtx, g.owner, mint)
	cancel()

	s.watch.mu.Lock()
	entry := s.watch.tokenAccounts[key]
	added := make([]string, 0)
	for _, account := range extra {
		if !containsString(entry.accounts, account) {
			entry.accounts = append(entry.accounts, account)
			added = append(added, account)
		}
	}
	// 只配置了一个节点时（如 publicnode 免费节点）该方法通常不可用，不记日志
	warn := listErr != nil && !entry.listErr && model.GetC(model.RpcEndpointSolanaPoll) != ""
	entry.listErr = listErr != nil
	entry.refreshAt = time.Now().Add(solanaTokenAccountsTTL)
	entry.refreshing = false
	s.watch.mu.Unlock()

	if warn {
		log.Task.Warn(fmt.Sprintf("Solana 查询钱包 %s 的代币账户失败，暂时只检测已知账户：%v", g.owner, listErr))
	}
	for _, account := range added {
		s.startPoll(ctx, g.target(account))
	}
}

// startPoll 在后台检测一个代币账户；该账户上一次检测还没结束时跳过
func (s *solana) startPoll(ctx context.Context, t solanaWatchTarget) {
	st := s.watch.account(t.account)

	s.watch.mu.Lock()
	if st.running {
		s.watch.mu.Unlock()

		return
	}
	st.running = true
	s.watch.inflight.Add(1)
	s.watch.mu.Unlock()

	go func() {
		defer s.watch.inflight.Done()
		defer func() {
			s.watch.mu.Lock()
			st.running = false
			s.watch.mu.Unlock()
		}()

		s.pollAccount(ctx, t)
	}()
}

func (s *solana) pollAccount(ctx context.Context, t solanaWatchTarget) {
	st := s.watch.account(t.account)
	if time.Now().Before(st.retryAt) {

		return
	}

	ctx, cancel := context.WithTimeout(ctx, solanaWatchPollTimeout)
	defer cancel()

	if !st.loaded {
		cursor, err := model.GetSolanaCursor(t.account)
		if err != nil {
			s.pollFailed(t, st, fmt.Errorf("read cursor: %w", err))

			return
		}
		pending, err := decodeSolanaPending(cursor.Pending)
		if err != nil {
			log.Task.Warn(fmt.Sprintf("Solana 代币账户 %s 的待处理签名无法解析，忽略：%v", t.account, err))
		}
		st.cursor = cursor
		st.pending = pending
		st.loaded = true
		log.Task.Info(fmt.Sprintf("Solana 开始按账户检测付款：%s（%s %s），低水位 slot %d，待处理签名 %d 笔", t.account, t.owner, t.tradeType, cursor.Slot, len(pending)))
	}

	items, complete, err := s.listSignatures(ctx, t.account, st.cursor.Slot, t.since)
	if err != nil {
		s.pollFailed(t, st, err)

		return
	}
	s.pollSucceeded(t, st, items)
	if !complete {
		log.Task.Error(fmt.Sprintf("Solana 代币账户 %s 新签名超过 %d 页，本轮只处理最新的部分，低水位不前进；较早的付款要等订单过期后对账", t.account, solanaWatchMaxPages))
	}

	// 列表里执行成功、还没处理完的签名都记入待处理表；执行失败的不必下载
	listed := make(map[string]bool, len(items))
	queue := make([]string, 0)
	for _, itm := range items {
		sig := itm.Get("signature").String()
		listed[sig] = true
		if !solanaTxSucceeded(itm.Get("err")) {
			s.watch.markDone(sig)

			continue
		}
		if s.watch.isDone(sig) {

			continue
		}
		st.pending[sig] = itm.Get("slot").Int()
		queue = append(queue, sig)
	}
	// 待处理表里这次没列出来的签名排在后面，按 slot 从新到旧，直接按签名下载
	unlisted := make([]string, 0)
	for sig := range st.pending {
		if !listed[sig] {
			unlisted = append(unlisted, sig)
		}
	}
	sort.Slice(unlisted, func(i, j int) bool { return st.pending[unlisted[i]] > st.pending[unlisted[j]] })
	queue = append(queue, unlisted...)

	// 从最新的签名开始处理，新付款不必排在旧签名后面；没下载过的签名排在重试的前面，
	// 一笔一直取不到（或节点卡住）的交易不会每轮都先占用时间，挡住同一账户后面的新付款
	now := time.Now()
	fresh := make([]string, 0, len(queue))
	retry := make([]string, 0)
	for _, sig := range queue {
		if !s.watch.due(sig, now) {

			continue
		}
		if s.watch.attempted(sig) {
			retry = append(retry, sig)

			continue
		}
		fresh = append(fresh, sig)
	}

	fetched := 0
	for _, sig := range append(fresh, retry...) {
		if fetched >= solanaWatchFetchPerRound || ctx.Err() != nil {

			break
		}

		fetched++
		fetchCtx, cancel := context.WithTimeout(ctx, solanaWatchFetchTimeout)
		err := s.processSignature(fetchCtx, sig)
		cancel()
		if err != nil {
			s.watch.retryLater(sig, time.Now(), err) // 退避从失败时算起，慢请求失败后不会立刻又到期

			continue
		}
		s.watch.markDone(sig)
	}

	for sig := range st.pending {
		if s.watch.isDone(sig) {
			delete(st.pending, sig)
		}
	}
	s.saveProgress(t, st, items, complete)
}

// listSignatures 从最新往旧列出代币账户的签名，到低水位 slot 或 since 为止；complete 为 false 表示翻页达到上限
func (s *solana) listSignatures(ctx context.Context, account string, floor int64, since time.Time) ([]gjson.Result, bool, error) {
	items := make([]gjson.Result, 0)
	limit := solanaWatchFirstPage
	opts := map[string]any{"limit": limit, "commitment": "confirmed"}
	for page := 0; page < solanaWatchMaxPages; page++ {
		result, err := s.rpcCall(ctx, "getSignaturesForAddress", account, opts)
		if err != nil {

			return nil, false, err
		}

		arr := result.Array()
		for _, itm := range arr {
			if itm.Get("slot").Int() < floor {

				return items, true, nil
			}
			if bt := itm.Get("blockTime"); bt.Type == gjson.Number && time.Unix(bt.Int(), 0).Before(since) {

				return items, true, nil
			}
			items = append(items, itm)
		}
		if len(arr) < limit {

			return items, true, nil
		}

		limit = solanaSignaturePageLimit
		opts = map[string]any{"limit": limit, "commitment": "confirmed", "before": arr[len(arr)-1].Get("signature").String()}
	}

	return items, false, nil
}

// processSignature 下载一笔交易并匹配订单；交易已被订单使用时直接返回
func (s *solana) processSignature(ctx context.Context, sig string) error {
	var used int64
	if err := model.Db.Model(&model.Order{}).Where("ref_hash = ?", sig).Count(&used).Error; err != nil {

		return err
	}
	if used > 0 {

		return nil
	}

	trans, err := s.rpcCallAccept(ctx, "getTransaction", solanaRequireTransaction, sig,
		map[string]any{"encoding": "json", "maxSupportedTransactionVersion": 1, "commitment": "confirmed"})
	if err != nil {

		return err
	}

	timestamp := time.Unix(trans.Get("blockTime").Int(), 0)

	return s.matchTransfers(s.parseTransaction(trans, int(trans.Get("slot").Int()), timestamp))
}

func solanaRequireTransaction(trans gjson.Result) error {
	if trans.Type == gjson.Null {

		return errSolanaNullResult
	}
	if trans.Get("blockTime").Type != gjson.Number || !trans.Get("meta").Exists() {

		return fmt.Errorf("transaction without blockTime or meta")
	}

	return nil
}

// matchTransfers 与 orderTransferHandle 的规则相同，但同步执行，并且用条件更新防止重复入账。
// MQTT 和"其它通知"需要整链数据，开启时由整链扫描负责推送（见 solanaMonitorEnabled），这里不重复推送
func (s *solana) matchTransfers(transfers []transfer) error {
	for _, t := range transfers {
		if !model.IsAmountValid(t.TradeType, t.Amount) {

			continue
		}

		if _, err := solanaMatchTransfer(t); err != nil {

			return err
		}
	}

	return nil
}

// solanaMatchTransfer 在同币种、同收款地址的订单中按创建时间找第一个匹配的订单，使其进入确认中。
// 候选包括过期不超过 solanaReconcileWindow 的订单：付款时间在有效期内、过期后才看到的付款照样入账，与订单过期后的对账一致
func solanaMatchTransfer(t transfer) (bool, error) {
	var orders []model.Order
	err := model.Db.Where("trade_type = ? and status in (?)", t.TradeType, []int{model.OrderStatusWaiting, model.OrderStatusExpired}).
		Where("(match_address = ? or (match_address = '' and address = ?))", t.RecvAddress, t.RecvAddress).
		Where("expired_at > ?", time.Now().Add(-solanaReconcileWindow)).
		Order("created_at asc").
		Find(&orders).Error
	if err != nil {

		return false, err
	}

	for _, o := range orders {
		if !orderTransferMatch(o, t) {

			continue
		}

		ok, err := solanaMarkConfirming(&o, t)
		if err != nil {

			return false, err
		}
		if !ok { // 订单状态刚变，或交易已被其他订单使用

			continue
		}

		log.Task.Info(fmt.Sprintf("Solana 识别付款 %s：%s %s %s，slot %d，延迟 %.1fs", o.TradeId, t.Amount, t.TradeType, t.TxHash, t.BlockNum, time.Since(t.Timestamp).Seconds()))

		return true, nil
	}

	return false, nil
}

// solanaMarkMu 让同一进程内各处（实时检测、对账）对 Solana 订单的入账判断串行执行，配合条件更新防止同一交易匹配两个订单
var solanaMarkMu sync.Mutex

func solanaMarkConfirming(o *model.Order, t transfer) (bool, error) {
	solanaMarkMu.Lock()
	defer solanaMarkMu.Unlock()

	return o.MarkConfirmingOnce(t.BlockNum, t.FromAddress, t.TxHash, t.Timestamp, t.Amount)
}

// saveProgress 把低水位推进到最早一笔未处理签名（列表里的和待处理表里的都算）；全部处理完时推进到最新签名之前 solanaWatchSlotOverlap 个 slot。
// 待处理表随低水位一起保存，重启后接着处理；翻页没到底（complete 为 false）时低水位不动，只保存待处理表
func (s *solana) saveProgress(t solanaWatchTarget, st *solanaAccountState, items []gjson.Result, complete bool) {
	low := st.cursor.Slot
	latest := ""
	if complete && len(items) > 0 {
		low = items[0].Get("slot").Int() - solanaWatchSlotOverlap
		for _, itm := range items {
			sig := itm.Get("signature").String()
			if !s.watch.isDone(sig) {
				low = min(low, itm.Get("slot").Int())

				continue
			}
			if latest == "" {
				latest = sig
			}
		}
		for _, slot := range st.pending {
			low = min(low, slot)
		}
	}
	low = max(low, st.cursor.Slot) // 低水位只前进；待处理签名都来自低水位之上的列表，不会被它越过

	pending := encodeSolanaPending(st.pending)
	if low == st.cursor.Slot && pending == st.cursor.Pending {

		return
	}

	cursor := model.SolanaCursor{Account: t.account, Owner: t.owner, TradeType: string(t.tradeType), Slot: low, Signature: latest, Pending: pending}
	if cursor.Signature == "" || low == st.cursor.Slot {
		cursor.Signature = st.cursor.Signature
	}
	if err := model.SaveSolanaCursor(cursor); err != nil {
		log.Task.Warn("Solana 保存检测进度失败：", err)

		return
	}
	st.cursor = cursor
}

// encodeSolanaPending 编码要保存的待处理签名；全部保存，否则重启后不知道低水位不能越过哪些签名
func encodeSolanaPending(pending map[string]int64) string {
	if len(pending) == 0 {

		return ""
	}

	data, _ := json.Marshal(pending) // 键按字母排序，内容不变时结果不变

	return string(data)
}

func decodeSolanaPending(raw string) (map[string]int64, error) {
	pending := make(map[string]int64)
	if raw == "" {

		return pending, nil
	}
	if err := json.Unmarshal([]byte(raw), &pending); err != nil {

		return make(map[string]int64), err
	}

	return pending, nil
}

func (s *solana) pollSucceeded(t solanaWatchTarget, st *solanaAccountState, items []gjson.Result) {
	slot := st.cursor.Slot
	if len(items) > 0 {
		slot = items[0].Get("slot").Int()
	}
	conf.RecordSuccess(conf.Solana, fmt.Sprint(slot))

	if st.alerted > 0 {
		log.Task.Info(fmt.Sprintf("Solana 代币账户 %s 恢复检测，此前连续失败 %s", t.account, time.Since(st.failingSince).Round(time.Second)))
	}
	st.failures = 0
	st.retryAt = time.Time{}
	st.failingSince = time.Time{}
	st.alerted = 0
}

func (s *solana) pollFailed(t solanaWatchTarget, st *solanaAccountState, err error) {
	conf.RecordFailure(conf.Solana)

	now := time.Now()
	st.failures++
	st.retryAt = now.Add(solanaWatchBackoff(st.failures))
	if st.failingSince.IsZero() {
		st.failingSince = now
	}

	failing := now.Sub(st.failingSince)
	switch {
	case failing >= solanaWatchErrorAfter && st.alerted < 2:
		st.alerted = 2
		log.Task.Error(fmt.Sprintf("Solana 代币账户 %s 已连续 %s 查询失败，新付款无法识别：%v", t.account, failing.Round(time.Second), err))
	case failing >= solanaWatchWarnAfter && st.alerted < 1:
		st.alerted = 1
		log.Task.Warn(fmt.Sprintf("Solana 代币账户 %s 已连续 %s 查询失败：%v", t.account, failing.Round(time.Second), err))
	default:
		log.Task.Warn(fmt.Sprintf("Solana 代币账户 %s 查询失败，%s 后重试：%v", t.account, solanaWatchBackoff(st.failures), err))
	}
}

// solanaWatchBackoff 返回第 n 次连续失败后的等待时间：2、4、8、16，最长 30 秒
func solanaWatchBackoff(n int) time.Duration {
	if n > 4 {

		return 30 * time.Second
	}

	return time.Duration(1<<n) * time.Second
}

// solanaSignatureBackoff 返回单笔交易第 n 次下载失败后的等待时间：2、4、8、16、32，最长 60 秒
func solanaSignatureBackoff(n int) time.Duration {
	if n > 5 {

		return time.Minute
	}

	return time.Duration(1<<n) * time.Second
}

func (w *solanaWatch) account(account string) *solanaAccountState {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.accounts[account]
	if !ok {
		st = &solanaAccountState{}
		w.accounts[account] = st
	}

	return st
}

// attempted 判断签名是否下载过（失败后在等重试）
func (w *solanaWatch) attempted(sig string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.signatures[sig]

	return ok && !st.done && st.attempts > 0
}

func (w *solanaWatch) isDone(sig string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.signatures[sig]

	return ok && st.done
}

// due 判断签名现在是否需要处理：没处理过，或者退避时间已到
func (w *solanaWatch) due(sig string, now time.Time) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.signatures[sig]

	return !ok || (!st.done && !now.Before(st.retryAt))
}

func (w *solanaWatch) markDone(sig string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.signatures[sig]
	if !ok {
		st = &solanaSignatureState{firstSeen: time.Now()}
		w.signatures[sig] = st
	}
	st.done = true
}

func (w *solanaWatch) retryLater(sig string, now time.Time, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	st, ok := w.signatures[sig]
	if !ok {
		st = &solanaSignatureState{firstSeen: now}
		w.signatures[sig] = st
	}
	st.attempts++
	if now.Sub(st.firstSeen) >= solanaWatchGiveUp {
		st.done = true
		log.Task.Error(fmt.Sprintf("Solana 交易 %s 连续 %s 无法下载，已放弃；如有订单，过期后由对账再查，需要人工核对：%v", sig, now.Sub(st.firstSeen).Round(time.Second), err))

		return
	}

	st.retryAt = now.Add(solanaSignatureBackoff(st.attempts))
	log.Task.Warn(fmt.Sprintf("Solana 交易 %s 下载失败（第 %d 次），%s 后单独重试：%v", sig, st.attempts, solanaSignatureBackoff(st.attempts), err))
}

func (w *solanaWatch) prune(now time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for sig, st := range w.signatures {
		if st.done && now.Sub(st.firstSeen) > solanaWatchSignatureTTL {
			delete(w.signatures, sig)
		}
	}
}

func containsString(arr []string, v string) bool {
	for _, itm := range arr {
		if itm == v {

			return true
		}
	}

	return false
}
