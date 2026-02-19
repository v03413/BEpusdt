package model

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/utils"
)

type OrderParams struct {
	Money         decimal.Decimal `json:"money"`          // 交易金额 (单位：法币)s
	ApiType       string          `json:"api_type"`       // 支付 API 类型
	Address       string          `json:"address"`        // 收款地址
	OrderId       string          `json:"order_id"`       // 商户订单 ID
	TradeType     TradeType       `json:"trade_type"`     // 交易类型
	RedirectUrl   string          `json:"redirect_url"`   // 成功跳转地址
	NotifyUrl     string          `json:"notify_url"`     // 异步通知地址
	Name          string          `json:"name"`           // 商品名称
	Timeout       int64           `json:"timeout"`        // 订单超时时间（秒）
	Rate          string          `json:"rate"`           // 强制指定汇率
	Fiat          Fiat            `json:"fiat"`           // 法币类型
	CurrencyLimit string          `json:"currency_limit"` // 限定币种
}

type Trade struct {
	Crypto  Crypto
	Rate    decimal.Decimal
	Address string
	Amount  string
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
	if _, ok := registry[p.TradeType]; !ok {

		return order, fmt.Errorf("不支持的交易类型：%s", p.TradeType)
	}
	if _, ok := supportFiat[p.Fiat]; !ok {

		return order, fmt.Errorf("不支持的法币类型：%s", p.Fiat)
	}

	maxAmount := decimal.NewFromFloat(cast.ToFloat64(GetC(PaymentMaxAmount)))
	minAmount := decimal.NewFromFloat(cast.ToFloat64(GetC(PaymentMinAmount)))
	if p.Money.GreaterThan(maxAmount) || p.Money.LessThan(minAmount) {

		return order, fmt.Errorf("交易金额必须在 %s - %s 之间", minAmount.String(), maxAmount.String())
	}

	Db.Where("order_id = ?", p.OrderId).Order("id desc").Limit(1).Find(&order)
	if order.Status == OrderStatusSuccess || order.Status == OrderStatusConfirming {
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
	if p.OrderId == t.OrderId &&
		p.TradeType == t.TradeType &&
		p.Money.String() == t.Money &&
		p.Fiat == t.Fiat {
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

	zero := time.Unix(0, 0)
	tradeOrder := Order{
		OrderId:       p.OrderId,
		TradeId:       tradeId,
		RefHash:       tradeId,
		TradeType:     p.TradeType,
		Rate:          fmt.Sprintf("%v", data.Rate),
		Amount:        data.Amount,
		Money:         p.Money.String(),
		Address:       data.Address,
		Status:        OrderStatusWaiting,
		Name:          p.Name,
		ApiType:       p.ApiType,
		ReturnUrl:     p.RedirectUrl,
		NotifyUrl:     p.NotifyUrl,
		NotifyNum:     0,
		NotifyState:   OrderNotifyStateFail,
		ExpiredAt:     CalcTradeExpiredAt(p.Timeout),
		Fiat:          p.Fiat,
		Crypto:        data.Crypto,
		CurrencyLimit: p.CurrencyLimit,
		ConfirmedAt:   &zero, // 默认填充一个0值时间，尽量避免数据库出现允许 NULL 值存在
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
	// 获取代币类型
	crypto, err := GetCrypto(p.TradeType)
	if err != nil {
		return Trade{}, fmt.Errorf("代币类型(%s)不支持：%v", p.TradeType, err)
	}

	// 获取订单汇率
	rate, err := GetOrderRate(crypto, p.Fiat, p.Rate)
	if err != nil {
		return Trade{}, err
	}
	if rate.LessThanOrEqual(decimal.Zero) {
		return Trade{}, fmt.Errorf("%s %s 汇率异常", crypto, p.Fiat)
	}

	var wallet = GetAvailableAddress(p.TradeType)
	if p.Address != "" {
		wallet = []string{p.Address}
		if !AddrCaseSens(p.TradeType) { // 交易类型不区分大小写，统一转小写；这个地址最后的交易匹配要用到，千万不能错
			wallet = []string{strings.ToLower(p.Address)}
		}
	}

	if len(wallet) == 0 {
		return Trade{}, fmt.Errorf("%s 未检测到可用钱包地址", p.TradeType)
	}

	// 计算交易金额
	address, amount, err := CalcTradeAmount(wallet, rate, p.Money, p.TradeType)
	if err != nil {

		return Trade{}, err
	}

	return Trade{
		Crypto:  crypto,
		Rate:    rate,
		Address: address,
		Amount:  amount,
	}, nil
}

// BuildPendingOrder 创建待支付订单（不锁定地址和汇率）
func BuildPendingOrder(p OrderParams) (Order, error) {
	var order Order

	maxAmount := decimal.NewFromFloat(cast.ToFloat64(GetC(PaymentMaxAmount)))
	minAmount := decimal.NewFromFloat(cast.ToFloat64(GetC(PaymentMinAmount)))
	if p.Money.GreaterThan(maxAmount) || p.Money.LessThan(minAmount) {
		return order, fmt.Errorf("交易金额必须在 %s - %s 之间", minAmount.String(), maxAmount.String())
	}

	Db.Where("order_id = ?", p.OrderId).Order("id desc").Limit(1).Find(&order)
	if order.Status == OrderStatusSuccess || order.Status == OrderStatusConfirming || order.Status == OrderStatusWaiting {
		return order, nil
	}

	// 默认使用 USDT
	crypto := USDT
	if p.TradeType != "" {
		c, err := GetCrypto(p.TradeType)
		if err == nil {
			crypto = c
		}
	} else if p.CurrencyLimit != "" {
		// 如果指定了唯一的币种限制，且不是黑名单模式，则使用该币种
		cur := strings.ToUpper(p.CurrencyLimit)
		if !strings.Contains(cur, ",") && !strings.HasPrefix(cur, "-") {
			crypto = Crypto(cur)
		}
	}

	return NewOrder(p, Trade{
		Crypto:  crypto,
		Rate:    decimal.Zero,
		Amount:  "0",
		Address: "",
	})
}
