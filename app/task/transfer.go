package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/smallnest/chanx"
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
	BlockNum    int64
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

func init() {
	Register(Task{Callback: orderTransferHandle})
	Register(Task{Callback: notOrderTransferHandle})
	Register(Task{Callback: tronResourceHandle})
}

func markFinalConfirmed(o model.Order) {
	o.SetSuccess()

	go notify.Handle(o)    // 订单回调
	go notifier.Success(o) // 消息通知
}

func orderTransferHandle(context.Context) {
	for transfers := range transferQueue.Out {
		var other = make([]transfer, 0)
		var orders = getAllWaitingOrders()
		for _, t := range transfers {
			// debug
			//if t.TradeType == model.UsdcBep20 {
			//	fmt.Println(t.TradeType, t.TxHash, t.FromAddress, "=>", t.RecvAddress, t.Amount.String())
			//}

			// 判断金额是否在允许范围内
			if !inAmountRange(t.Amount) {

				continue
			}

			// 判断是否存在对应订单
			o, ok := orders[fmt.Sprintf("%s%v%s", t.RecvAddress, t.Amount.String(), t.TradeType)]
			if !ok {
				other = append(other, t)

				continue
			}

			// 有效期检测
			if !o.CreatedAt.Before(t.Timestamp) || !o.ExpiredAt.After(t.Timestamp) {

				continue
			}

			// 进入确认状态
			o.MarkConfirming(t.BlockNum, t.FromAddress, t.TxHash, t.Timestamp)
		}

		if len(other) > 0 {
			notOrderQueue.In <- other
		}
	}
}

func notOrderTransferHandle(context.Context) {
	for transfers := range notOrderQueue.Out {
		var was []model.Wallet

		model.Db.Where("other_notify = ?", model.WaOtherEnable).Find(&was)

		for _, wa := range was {
			for _, t := range transfers {
				if t.RecvAddress != wa.Address && t.FromAddress != wa.Address {

					continue
				}

				if !inAmountRange(t.Amount) {

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
	}
}

func tronResourceHandle(context.Context) {
	for resources := range resourceQueue.Out {
		var was []model.Wallet
		var types = []model.TradeType{model.TronTrx, model.UsdtTrc20}

		model.Db.Where("status = ? and other_notify = ? and trade_type in (?)", model.WaStatusEnable, model.WaOtherEnable, types).Find(&was)

		for _, wa := range was {
			for _, t := range resources {
				if t.RecvAddress != wa.Address && t.FromAddress != wa.Address {

					continue
				}

				if t.ResourceCode != core.ResourceCode_ENERGY {

					continue
				}

				if !model.IsNeedNotifyByTxid(t.ID) {

					continue
				}

				var record = model.NotifyRecord{Txid: t.ID}
				model.Db.Create(&record)

				notifier.TronResourceChange(model.TronResource(t))
			}
		}
	}
}

func getAllWaitingOrders() map[string]model.Order {
	var tradeOrders = model.GetOrderByStatus(model.OrderStatusWaiting)
	var data = make(map[string]model.Order) // 当前所有正在等待支付的订单 Lock Key
	for _, order := range tradeOrders {
		if time.Now().Unix() >= order.ExpiredAt.Unix() { // 订单过期
			order.SetExpired()
			notify.Bepusdt(order)

			continue
		}

		if order.TradeType == model.UsdtPolygon {

			order.Address = strings.ToLower(order.Address)
		}

		data[order.Address+order.Amount+string(order.TradeType)] = order
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

func inAmountRange(payAmount decimal.Decimal) bool {
	var payMin, payMax = model.Payment()

	if payAmount.GreaterThan(payMax) {

		return false
	}

	if payAmount.LessThan(payMin) {

		return false
	}

	return true
}
