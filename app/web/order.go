package web

import (
	"fmt"
	"sync"
	"time"

	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/help"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/task/rate"
)

type orderParams struct {
	Money       float64 `json:"money"`        // 交易金额 CNY
	ApiType     string  `json:"api_type"`     // 支付API类型
	PayAddress  string  `json:"pay_address"`  // 收款地址
	OrderId     string  `json:"order_id"`     // 商户订单ID
	TradeType   string  `json:"trade_type"`   // 交易类型
	RedirectUrl string  `json:"redirect_url"` // 成功跳转地址
	NotifyUrl   string  `json:"notify_url"`   // 异步通知地址
	Name        string  `json:"name"`         // 商品名称
	Timeout     uint64  `json:"timeout"`      // 订单超时时间（秒）
	Rate        string  `json:"rate"`         // 强制指定汇率
}

func buildOrder(p orderParams) (model.TradeOrders, error) {
	var lock sync.Mutex
	var order model.TradeOrders

	model.DB.Where("order_id = ?", p.OrderId).Find(&order)
	if order.Status == model.OrderStatusWaiting || order.Status == model.OrderStatusSuccess {
		return order, nil
	}

	// 暂时先强制使用互斥锁，后续有需求的话再考虑优化
	lock.Lock()
	defer lock.Unlock()

	// 获取兑换汇率
	var calcRate float64
	if p.Rate != "" {
		rawRate := rate.GetOkxUsdtRawRate()
		if p.TradeType == model.OrderTradeTypeTronTrx {
			rawRate = rate.GetOkxTrxRawRate()
		}
		calcRate = rate.ParseFloatRate(p.Rate, rawRate)
	} else {
		calcRate = rate.GetUsdtCalcRate()
		if p.TradeType == model.OrderTradeTypeTronTrx {
			calcRate = rate.GetTrxCalcRate()
		}
	}

	// 获取钱包地址
	wallet := model.GetAvailableAddress(p.PayAddress, p.TradeType)
	if len(wallet) == 0 {
		return order, fmt.Errorf(fmt.Sprintf("类型(%s)没检测到可用收款地址", p.TradeType))
	}

	// 计算交易金额
	address, amount := model.CalcTradeAmount(wallet, calcRate, p.Money, p.TradeType)
	tradeId, err := help.GenerateTradeId()
	if err != nil {
		return order, err
	}

	// 超时处理
	timeout := conf.GetExpireTime() * time.Second
	if p.Timeout >= 60 { // 至少60秒

		timeout = time.Duration(p.Timeout) * time.Second
	}

	// 创建交易订单
	expiredAt := time.Now().Add(timeout)
	tradeOrder := model.TradeOrders{
		OrderId:     p.OrderId,
		TradeId:     tradeId,
		TradeHash:   tradeId, // 这里默认填充一个本地交易ID，等支付成功后再更新为实际交易哈希
		TradeType:   p.TradeType,
		TradeRate:   fmt.Sprintf("%v", calcRate),
		Amount:      amount,
		Money:       p.Money,
		Address:     address.Address,
		Status:      model.OrderStatusWaiting,
		Name:        p.Name,
		ApiType:     p.ApiType,
		ReturnUrl:   p.RedirectUrl,
		NotifyUrl:   p.NotifyUrl,
		NotifyNum:   0,
		NotifyState: model.OrderNotifyStateFail,
		ExpiredAt:   expiredAt,
	}
	if err = model.DB.Create(&tradeOrder).Error; err != nil {
		log.Error("订单创建失败：", err.Error())

		return order, err
	}

	model.PushWebhookEvent(model.WebhookEventOrderCreate, tradeOrder)

	return tradeOrder, nil
}
