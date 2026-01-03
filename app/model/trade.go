package model

import (
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/utils"
)

type OrderParams struct {
	Money       decimal.Decimal `json:"money"`        // 交易金额 (单位：法币)s
	ApiType     string          `json:"api_type"`     // 支付 API 类型
	Address     string          `json:"address"`      // 收款地址
	OrderId     string          `json:"order_id"`     // 商户订单 ID
	TradeType   TradeType       `json:"trade_type"`   // 交易类型
	RedirectUrl string          `json:"redirect_url"` // 成功跳转地址
	NotifyUrl   string          `json:"notify_url"`   // 异步通知地址
	Name        string          `json:"name"`         // 商品名称
	Timeout     int64           `json:"timeout"`      // 订单超时时间（秒）
	Rate        string          `json:"rate"`         // 强制指定汇率
	Fiat        string          `json:"fiat"`         // 法币类型
}

type Trade struct {
	TokenType TokenType
	Rate      decimal.Decimal
	Address   string
	Amount    string
}

func BuildOrder(p OrderParams) (Order, error) {
	var order Order

	if p.Address != "" {
		if !utils.IsValidTronAddress(p.Address) &&
			!utils.IsValidEvmAddress(p.Address) &&
			!utils.IsValidSolanaAddress(p.Address) &&
			!utils.IsValidAptosAddress(p.Address) {

			return order, fmt.Errorf("钱包地址格式错误：%s", p.Address)
		}
	}
	if _, ok := SupportTradeTypes[p.TradeType]; !ok {

		return order, fmt.Errorf("不支持的交易类型：%s", p.TradeType)
	}
	if !utils.InStrings(p.Fiat, SupportTradeFiat) {

		return order, fmt.Errorf("不支持的法币类型：%s", p.Fiat)
	}

	Db.Where("order_id = ?", p.OrderId).Find(&order)
	if order.Status == OrderStatusSuccess {
		return order, nil
	}
	if order.Status == OrderStatusWaiting {
		return RebuildOrder(order, p)
	}

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	data, err := BuildTrade(p)
	if err != nil {
		return order, err
	}

	return NewOrder(p, data)
}

func RebuildOrder(t Order, p OrderParams) (Order, error) {
	if p.OrderId == t.OrderId && p.TradeType == t.TradeType && p.Money.String() == t.Money {
		return t, nil
	}

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	data, err := BuildTrade(p)
	if err != nil {
		return t, err
	}

	t.Amount = data.Amount
	t.TradeType = p.TradeType
	t.Address = data.Address
	t.ExpiredAt = CalcTradeExpiredAt(p.Timeout)

	return t, Db.Save(&t).Error
}

func NewOrder(p OrderParams, data Trade) (Order, error) {
	tradeId, err := utils.GenerateTradeId()
	if err != nil {
		return Order{}, err
	}

	tradeOrder := Order{
		OrderId:     p.OrderId,
		TradeId:     tradeId,
		RefHash:     tradeId,
		TradeType:   p.TradeType,
		Rate:        fmt.Sprintf("%v", data.Rate),
		Amount:      data.Amount,
		Money:       p.Money.String(),
		Address:     data.Address,
		Status:      OrderStatusWaiting,
		Name:        p.Name,
		ApiType:     p.ApiType,
		ReturnUrl:   p.RedirectUrl,
		NotifyUrl:   p.NotifyUrl,
		NotifyNum:   0,
		NotifyState: OrderNotifyStateFail,
		ExpiredAt:   CalcTradeExpiredAt(p.Timeout),
		Fiat:        p.Fiat,
	}

	if tradeOrder.Name == "" {
		tradeOrder.Name = tradeOrder.OrderId
	}

	if err = Db.Create(&tradeOrder).Error; err != nil {
		log.Error("订单创建失败：", err.Error())
		return Order{}, err
	}

	return tradeOrder, nil
}

func BuildTrade(p OrderParams) (Trade, error) {
	var t = TradeType(p.TradeType)

	// 获取代币类型
	tokenType, err := GetTokenType(t)
	if err != nil {
		return Trade{}, fmt.Errorf("代币类型(%s)不支持：%v", p.TradeType, err)
	}

	// 获取订单汇率
	rate, err := getOrderRate(tokenType, p.Fiat, p.Rate)
	if err != nil {
		return Trade{}, err
	}
	if rate.LessThanOrEqual(decimal.Zero) {
		return Trade{}, fmt.Errorf("%s %s 汇率异常", tokenType, p.Fiat)
	}

	var wallet = GetAvailableAddress(p.TradeType)
	if p.Address != "" {
		wallet = []string{p.Address}
	}

	if len(wallet) == 0 {
		return Trade{}, fmt.Errorf("%s 未检测到可用钱包地址", p.TradeType)
	}

	// 计算交易金额
	address, amount := CalcTradeAmount(wallet, rate, p.Money, t)

	return Trade{
		TokenType: tokenType,
		Rate:      rate,
		Address:   address,
		Amount:    amount,
	}, nil
}
