package task

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/shopspring/decimal"
	"github.com/smallnest/chanx"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	tgo "github.com/xssnick/tonutils-go/ton"
)

const tonMasterChainID = -1
const tonMasterShard = int64(-0x8000000000000000)
const tonOpInternalTransfer = uint32(0x178d4519)

type ton struct {
	lastBlockSeqno uint32
	shardTipMap    sync.Map // key: "workchain:shard", value: uint32 last processed shard seqno
	blockScanQueue *chanx.UnboundedChan[uint32]
	clientOnce     sync.Once
	api            tgo.APIClientWrapped
}

var tn ton

func init() {
	tn = ton{
		blockScanQueue: chanx.NewUnboundedChan[uint32](context.Background(), 30),
	}

	Register(Task{Callback: tn.syncMBSeqnoForward})
	Register(Task{Duration: time.Second, Callback: tn.blockDispatch})
	Register(Task{Duration: time.Second * 3, Callback: tn.tradeConfirmHandle})
}

func (t *ton) syncMBSeqnoForward(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if t.syncBreak() {
			time.Sleep(time.Second * 3)
			continue
		}

		// 初始化：首次获取当前最新高度作为起点
		if t.lastBlockSeqno == 0 {
			mb, err := t.client().CurrentMasterchainInfo(ctx)
			if err != nil {
				log.Task.Warn(fmt.Sprintf("get current masterchain info: %v", err))
				time.Sleep(time.Second)
				continue
			}

			t.syncMBSeqnoBackward(mb.SeqNo)
			t.lastBlockSeqno = mb.SeqNo - 1
		}

		nextSeqno := t.lastBlockSeqno + 1

		// 阻塞直到 nextSeqno 在节点上可用，无需轮询
		mb, err := t.client().WaitForBlock(nextSeqno).GetMasterchainInfo(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			log.Task.Warn(fmt.Sprintf("WaitForBlock GetMasterchainInfo err: %v", err))
			time.Sleep(time.Second)
			continue
		}

		now := mb.SeqNo

		// 区块高度变化过大，强制丢块重扫
		if now-t.lastBlockSeqno > cast.ToUint32(model.GetC(model.BlockHeightMaxDiff)) {
			t.lastBlockSeqno = now - 1
		}

		// 待扫描区块入列
		for n := t.lastBlockSeqno + 1; n <= now; n++ {
			t.blockScanQueue.In <- n
		}

		t.lastBlockSeqno = now
	}
}

func (t *ton) syncMBSeqnoBackward(now uint32) {
	if now == 0 || t.lastBlockSeqno != 0 {

		return
	}

	var o model.Order
	trade := model.GetNetworkTrades(conf.Ton)
	model.Db.Model(&model.Order{}).Where("status = ? and trade_type in (?)", model.OrderStatusWaiting, trade).Order("created_at asc").Limit(1).Find(&o)
	if o.ID == 0 {

		return
	}

	// 大概1秒3个区块（大概值，实际存在波动）
	num := uint32((time.Now().Unix() - o.CreatedAt.Time().Unix() + 1) * 3) // 计算需要反向扫描的区块数量

	go func() {
		ticker := time.NewTicker(time.Millisecond * 125)
		defer ticker.Stop()

		var i uint32

		for i = 0; i < num; i++ {
			if t.syncBreak() {

				return
			}

			t.blockScanQueue.In <- now - i

			<-ticker.C
		}
	}()
}

func (t *ton) blockDispatch(ctx context.Context) {
	p, err := ants.NewPoolWithFunc(3, t.blockParse)
	if err != nil {
		log.Task.Warn("Error creating pool:", err)

		return
	}

	defer p.Release()

	for {
		select {
		case <-ctx.Done():
			return
		case n, ok := <-t.blockScanQueue.Out:
			if !ok {
				return
			}
			if err := p.Invoke(n); err != nil {

				log.Task.Warn("Tron Error invoking process block:", err)
			}
		}
	}
}

func (t *ton) blockParse(n any) {
	var seqno = n.(uint32)
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	mb, err := t.client().LookupBlock(ctx, tonMasterChainID, tonMasterShard, seqno)
	if err != nil {
		conf.RecordFailure(conf.Ton)
		log.Task.Warn("Ton LookupBlock ", err)

		return
	}

	conf.RecordSuccess(conf.Ton, cast.ToString(seqno))

	// 目前实际来看，basechain 还未分裂，所以 len(shardsTip) == 1
	shardsTip, err := t.client().GetBlockShardsInfo(ctx, mb)
	if err != nil {
		conf.RecordFailure(conf.Ton)
		log.Task.Warn(fmt.Sprintf("get shards info seqno=%d err: %v", seqno, err))

		return
	}

	for _, tip := range shardsTip {
		start := tip.SeqNo
		key := fmt.Sprintf("%d:%d", tip.Workchain, tip.Shard)
		if old, loaded := t.shardTipMap.Swap(key, tip.SeqNo); loaded {
			start = old.(uint32) + 1
		}

		for s := start; s <= tip.SeqNo; s++ {
			shardBlock, err := t.client().LookupBlock(ctx, tip.Workchain, tip.Shard, s)
			if err != nil {
				conf.RecordFailure(conf.Ton)
				log.Task.Warn(fmt.Sprintf("lookup shard block workchain=%d shard=%d seqno=%d err: %v", tip.Workchain, tip.Shard, s, err))

				continue
			}

			if err := t.processShard(ctx, shardBlock, seqno); err != nil {
				conf.RecordFailure(conf.Ton)
				log.Task.Warn(fmt.Sprintf("process shard block workchain=%d shard=%d seqno=%d err: %v", tip.Workchain, tip.Shard, s, err))
			}
		}
	}

	log.Task.Info(fmt.Sprintf("区块扫描完成(Ton): %d 成功率：%s", seqno, conf.GetSuccessRate(conf.Ton)))
}

func (t *ton) syncBreak() bool {
	if t.blockScanQueue.Len() >= blockQueueLimit {
		log.Task.Warn("ton 同步阻塞，当前区块消费堆积数量：", t.blockScanQueue.Len())

		return true
	}

	if mqttSubscribed(conf.Ton) {
		return false
	}

	var count int64 = 0
	trade := []model.TradeType{model.UsdtTon, model.TonGram}
	model.Db.Model(&model.Order{}).Where("status = ? and trade_type in (?)", model.OrderStatusWaiting, trade).Count(&count)
	if count > 0 {

		return false
	}

	model.Db.Model(&model.Wallet{}).Where("other_notify = ? and trade_type in (?)", model.WaOtherEnable, trade).Count(&count)
	if count > 0 {

		return false
	}

	return true
}

func (t *ton) processShard(ctx context.Context, shard *tgo.BlockIDExt, seqno uint32) error {
	data, err := t.client().GetBlockData(ctx, shard)
	if err != nil {
		return err
	}
	txs, err := data.ListTransactions()
	if err != nil {
		return err
	}

	var transfers = make([]transfer, 0)
	for _, tx := range txs {
		if tonTrans, ok := t.parseTonTransfer(tx, seqno); ok {
			transfers = append(transfers, tonTrans)
		}
		if jettonTrans, ok := t.parseInternalTransfer(shard, tx, seqno); ok {
			transfers = append(transfers, jettonTrans)
		}
	}

	if len(transfers) > 0 {
		transferQueue.In <- transfers
	}

	return nil
}

// https://github.com/ton-blockchain/TEPs/blob/63fc78718dd9930f3e106954ebec743c3ad07993/text/0074-jettons-standard.md?plain=1#L226
// https://github.com/ton-blockchain/token-contract/blob/1182ad99413242f09925d50e70ccb7e0e09f94d4/ft/jetton-wallet.fc#L43
func (t *ton) parseInternalTransfer(shard *tgo.BlockIDExt, tx *tlb.Transaction, blockNum uint32) (transfer, bool) {
	if tx.IO.In == nil {
		return transfer{}, false
	}
	in, ok := tx.IO.In.Msg.(*tlb.InternalMessage)
	if !ok || in.Bounced || in.Body == nil {
		return transfer{}, false
	}
	s, err := in.Body.BeginParse()
	if err != nil || s.BitsLeft() < 32 {
		return transfer{}, false
	}
	v, err := s.LoadUInt(32)
	if err != nil || uint32(v) != tonOpInternalTransfer {
		return transfer{}, false
	}

	transOrd, ok := tx.Description.(tlb.TransactionDescriptionOrdinary)
	if !ok {
		return transfer{}, false
	}

	compute, ok := transOrd.ComputePhase.Phase.(tlb.ComputePhaseVM)
	if !ok || !compute.Success {
		return transfer{}, false
	}

	if transOrd.ActionPhase == nil || !transOrd.ActionPhase.Success {
		return transfer{}, false
	}

	if err = s.SkipBits(64); err != nil {
		return transfer{}, false
	}
	amount, err := s.LoadBigCoins()
	if err != nil || amount.Sign() <= 0 {
		return transfer{}, false
	}
	fromOwner, err := s.LoadAddr()
	if err != nil {
		return transfer{}, false
	}
	if fromOwner.Type() != address.StdAddress {
		fromOwner = address.NewAddress(0, byte(fromOwner.Workchain()), fromOwner.Data())
	}

	toJetton := address.NewAddress(0, byte(shard.Workchain), tx.AccountAddr)

	return transfer{
		Network:     conf.Ton,
		TxHash:      hex.EncodeToString(tx.Hash),
		Amount:      decimal.NewFromBigInt(amount, conf.UsdtTonDecimals),
		FromAddress: fromOwner.Bounce(false).String(),
		RecvAddress: toJetton.Bounce(false).String(),
		Timestamp:   time.Unix(int64(tx.Now), 0),
		TradeType:   model.UsdtTon,
		BlockNum:    int(blockNum),
	}, true
}

func (t *ton) parseTonTransfer(tx *tlb.Transaction, blockNum uint32) (transfer, bool) {
	in := tx.IO.In
	if in == nil {
		return transfer{}, false
	}

	msg, ok := in.Msg.(*tlb.InternalMessage)
	if !ok {
		return transfer{}, false
	}
	if msg.Bounced {
		return transfer{}, false
	}
	if msg.Amount.Nano().Sign() <= 0 {
		return transfer{}, false
	}
	if msg.Body.BitsSize() != 0 {
		return transfer{}, false
	}

	return transfer{
		Network:     conf.Ton,
		TxHash:      hex.EncodeToString(tx.Hash),
		FromAddress: msg.SrcAddr.Bounce(false).String(),
		RecvAddress: msg.DstAddr.Bounce(false).String(),
		Timestamp:   time.Unix(int64(tx.Now), 0),
		Amount:      decimal.NewFromBigInt(msg.Amount.Nano(), conf.TonTonDecimals),
		TradeType:   model.TonGram,
		BlockNum:    int(blockNum),
	}, true
}

func (t *ton) tradeConfirmHandle(context.Context) {
	var orders = getConfirmingOrders([]model.TradeType{model.UsdtTon, model.TonGram})
	var wg sync.WaitGroup

	for _, order := range orders {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// 一旦某笔交易所在的 shard block 被 MasterChain block 引用（commit），则该交易获得最终性（finality）。
			markFinalConfirmed(order)
		}()
	}

	wg.Wait()
}

func (t *ton) client() tgo.APIClientWrapped {
	t.clientOnce.Do(func() {
		t.api = utils.NewTonClient(model.GetC(model.RpcGlobalConfigUrlTon))
	})

	return t.api
}
