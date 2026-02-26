package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/notifier"
	"github.com/v03413/bepusdt/app/task/notify"
	"github.com/v03413/tronprotocol/core"
)

type transfer struct {
	Network     string
	TxHash      string
	Amount      decimal.Decimal
	FromAddress string
	RecvAddress string
	Timestamp   time.Time
	TradeType   model.TradeType
	BlockNum    int
}

type resource struct {
	ID           string
	Type         core.Transaction_Contract_ContractType
	Balance      int64
	FromAddress  string
	RecvAddress  string
	Timestamp    time.Time
	ResourceCode core.ResourceCode
}

var resourceQueue = chanx.NewUnboundedChan[[]resource](context.Background(), 30) // 资源队列
var notOrderQueue = chanx.NewUnboundedChan[[]transfer](context.Background(), 30) // 非订单队列
var transferQueue = chanx.NewUnboundedChan[[]transfer](context.Background(), 30) // 交易转账队列

const batchInterval = time.Second * 1       // 批处理缓解数据库读取压力
const orderCheckInterval = time.Second * 10 // 订单过期检查间隔

func init() {
	Register(Task{Callback: orderTransferHandle})
	Register(Task{Callback: notOrderTransferHandle})
	Register(Task{Callback: tronResourceHandle})
}

func orderTransferHandle(ctx context.Context) {
	var batch = make([]transfer, 0, 1000)
	var lastCheckTime = time.Now()
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case transfers, ok := <-transferQueue.Out:
			if !ok {
				return
			}
			batch = append(batch, transfers...)
		case <-ticker.C:
			// 每10秒强制检查一次过期订单，即使没有交易，防止无交易时订单不过期
			var shouldCheck = time.Since(lastCheckTime) >= orderCheckInterval
			if len(batch) == 0 && !shouldCheck {
				continue
			}

			if shouldCheck {
				lastCheckTime = time.Now()
			}

			var other = make([]transfer, 0)
			// getAllWaitingOrders 内部包含过期检查逻辑
			var orders = getAllWaitingOrders()
			if len(batch) == 0 {
				continue
			}
			for _, t := range batch {
				// 判断数额是否在允许范围内
				if !model.IsAmountValid(t.TradeType, t.Amount) {
					continue
				}

				key := fmt.Sprintf("%s%s", t.RecvAddress, t.TradeType)
				orderList, ok := orders[key]
				if !ok {
					other = append(other, t)
					continue
				}

				var matched bool
				for _, o := range orderList {
					// 金额前缀匹配
					if !amountMatch(t.Amount, o.Amount, string(o.TradeType)) {

						continue
					}

					// 有效期检测
					if !o.CreatedAt.Before(t.Timestamp) || !o.ExpiredAt.After(t.Timestamp) {

						continue
					}

					// 找到匹配的订单
					o.MarkConfirming(t.BlockNum, t.FromAddress, t.TxHash, t.Timestamp)
					matched = true
					break
				}

				if !matched {
					other = append(other, t)
				}
			}

			if len(other) > 0 {
				notOrderQueue.In <- other
			}

			batch = batch[:0]
		}
	}
}

func notOrderTransferHandle(ctx context.Context) {
	var batch = make([]transfer, 0, 1000)
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case transfers, ok := <-notOrderQueue.Out:
			if !ok {
				return
			}
			batch = append(batch, transfers...)
		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}

			var was = make([]model.Wallet, 0)
			model.Db.Where("other_notify = ?", model.WaOtherEnable).Find(&was)
			for _, wa := range was {
				for _, t := range batch {
					if t.RecvAddress != wa.MatchAddr && t.FromAddress != wa.MatchAddr {
						continue
					}

					if !model.IsNeedNotifyByTxid(t.TxHash) {
						continue
					}

					var record = model.NotifyRecord{Txid: t.TxHash}
					model.Db.Create(&record)

					notifier.NonOrderTransfer(model.TronTransfer(t), wa)
				}
			}

			batch = batch[:0]
		}
	}
}

func tronResourceHandle(ctx context.Context) {
	var batch = make([]resource, 0, 1000)
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case resources, ok := <-resourceQueue.Out:
			if !ok {
				return
			}
			batch = append(batch, resources...)
		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}

			var was []model.Wallet
			model.Db.Where("status = ? and other_notify = ?", model.WaStatusEnable, model.WaOtherEnable).Find(&was)

			for _, wa := range was {
				if wa.GetNetwork() != conf.Tron {
					// 只有 Tron 网络目前才有资源变更通知
					continue
				}

				for _, r := range batch {
					if r.RecvAddress != wa.Address && r.FromAddress != wa.Address {
						continue
					}
					if r.ResourceCode != core.ResourceCode_ENERGY {
						continue
					}
					if !model.IsNeedNotifyByTxid(r.ID) {
						continue
					}

					var record = model.NotifyRecord{Txid: r.ID}
					model.Db.Create(&record)

					notifier.TronResourceChange(model.TronResource(r))
				}
			}

			batch = batch[:0]
		}
	}
}

func markFinalConfirmed(o model.Order) {
	o.SetSuccess()

	go notify.Handle(o)    // 订单回调
	go notifier.Success(o) // 消息通知
}

func getAllWaitingOrders() map[string][]model.Order {
	var tradeOrders = model.GetOrderByStatus(model.OrderStatusWaiting)
	var data = make(map[string][]model.Order)
	for _, t := range tradeOrders {
		if time.Now().Unix() >= t.ExpiredAt.Unix() { // 订单过期
			t.SetExpired()
			notify.Bepusdt(t)

			continue
		}

		key := t.Address + string(t.TradeType)
		data[key] = append(data[key], t)
	}

	return data
}

func getConfirmingOrders(tradeType []model.TradeType) []model.Order {
	var orders = make([]model.Order, 0)
	var data = make([]model.Order, 0)
	var db = model.Db.Where("status = ?", model.OrderStatusConfirming)
	if len(tradeType) > 0 {
		db = db.Where("trade_type in (?)", tradeType)
	}

	db.Find(&orders)

	for _, order := range orders {
		if time.Now().Unix() >= order.ExpiredAt.Unix() {
			order.SetFailed()
			notify.Bepusdt(order)

			continue
		}

		data = append(data, order)
	}

	return data
}

func amountMatch(amount decimal.Decimal, target, tradeType string) bool {
	mode := model.GetC(model.PaymentMatchMode)
	switch model.MatchMode(mode) {
	case model.Classic:
		return amount.String() == target
	case model.HasPrefix:
		return strings.HasPrefix(amount.String(), target)
	case model.RoundOff:
		t, err := decimal.NewFromString(target)
		if err != nil {

			return false
		}

		_, precision := model.GetAtomicity(model.TradeType(tradeType)) // 标准精度
		precision2 := abs(t.Exponent())                                // 实际精度
		if precision2 != precision {

			precision = precision2
		}

		a := amount.Round(precision)
		t = t.Round(precision)

		return a.Equal(t)
	}

	return false
}

func abs(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}
