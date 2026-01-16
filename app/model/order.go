package model

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

const (
	OrderNotifyStateSucc = 1 // 回调成功
	OrderNotifyStateFail = 0 // 回调失败

	OrderStatusWaiting    = 1 // 等待支付
	OrderStatusSuccess    = 2 // 交易确认成功
	OrderStatusExpired    = 3 // 订单过期
	OrderStatusCanceled   = 4 // 订单取消
	OrderStatusConfirming = 5 // 等待交易确认
	OrderStatusFailed     = 6 // 交易确认失败

	BscBnb      TradeType = "bsc.bnb"
	EthereumEth TradeType = "ethereum.eth"
	TronTrx     TradeType = "tron.trx"

	UsdtTrc20    TradeType = "usdt.trc20"
	UsdcTrc20    TradeType = "usdc.trc20"
	UsdtPolygon  TradeType = "usdt.polygon"
	UsdcPolygon  TradeType = "usdc.polygon"
	UsdtArbitrum TradeType = "usdt.arbitrum"
	UsdcArbitrum TradeType = "usdc.arbitrum"
	UsdtErc20    TradeType = "usdt.erc20"
	UsdcErc20    TradeType = "usdc.erc20"
	UsdtBep20    TradeType = "usdt.bep20"
	UsdcBep20    TradeType = "usdc.bep20"
	UsdtXlayer   TradeType = "usdt.xlayer"
	UsdcXlayer   TradeType = "usdc.xlayer"
	UsdcBase     TradeType = "usdc.base"
	UsdtSolana   TradeType = "usdt.solana"
	UsdcSolana   TradeType = "usdc.solana"
	UsdtAptos    TradeType = "usdt.aptos"
	UsdcAptos    TradeType = "usdc.aptos"
	UsdtPlasma   TradeType = "usdt.plasma"
)

const (
	OrderApiTypeEpusdt = "epusdt" // epusdt
	OrderApiTypeEpay   = "epay"   // 彩虹易支付
)

var calcMutex sync.Mutex

type Order struct {
	Id
	OrderId     string     `gorm:"column:order_id;type:varchar(128);not null;index;comment:商户ID" json:"order_id"`
	TradeId     string     `gorm:"column:trade_id;type:varchar(128);not null;uniqueIndex;comment:本地ID" json:"trade_id"`
	TradeType   TradeType  `gorm:"column:trade_type;type:varchar(20);not null;index;comment:交易类型" json:"trade_type"`
	Fiat        Fiat       `gorm:"column:fiat;type:varchar(16);not null;index;default:CNY;comment:法定货币" json:"fiat"`
	Crypto      Crypto     `gorm:"column:crypto;type:varchar(16);not null;index;default:USDT;comment:加密货币" json:"crypto"`
	Rate        string     `gorm:"column:rate;type:varchar(10);not null;comment:交易汇率" json:"rate"`
	Amount      string     `gorm:"column:amount;type:varchar(32);not null;default:0.00;comment:交易数额" json:"amount"`
	Money       string     `gorm:"column:money;type:varchar(32);not null;default:0.00;comment:交易金额" json:"money"`
	Address     string     `gorm:"column:address;type:varchar(128);index;not null;comment:收款地址" json:"address"`
	FromAddress string     `gorm:"column:from_address;type:varchar(128);not null;default:'';comment:支付地址" json:"from_address"`
	Status      int        `gorm:"column:status;type:tinyint(1);not null;default:1;index;comment:交易状态" json:"status"`
	Name        string     `gorm:"column:name;type:varchar(64);not null;default:'';comment:商品名称" json:"name"`
	ApiType     string     `gorm:"column:api_type;type:varchar(20);not null;default:'epusdt';comment:API类型" json:"api_type"`
	ReturnUrl   string     `gorm:"column:return_url;type:varchar(255);not null;default:'';comment:同步地址" json:"return_url"`
	NotifyUrl   string     `gorm:"column:notify_url;type:varchar(255);not null;default:'';comment:异步地址" json:"notify_url"`
	NotifyNum   int        `gorm:"column:notify_num;type:int(11);not null;default:0;comment:回调次数" json:"notify_num"`
	NotifyState int        `gorm:"column:notify_state;type:tinyint(1);not null;default:0;comment:回调状态 1：成功 0：失败" json:"notify_state"`
	RefHash     string     `gorm:"column:ref_hash;type:varchar(128);default:'';index;comment:交易哈希" json:"ref_hash"`
	RefBlockNum int64      `gorm:"column:ref_block_num;type:bigint(20);not null;default:0;comment:区块索引" json:"ref_block_num"`
	ExpiredAt   time.Time  `gorm:"column:expired_at;type:timestamp;not null;comment:失效时间" json:"expired_at"`
	ConfirmedAt *time.Time `gorm:"column:confirmed_at;type:timestamp;null;comment:交易确认时间" json:"confirmed_at"`
	AutoTimeAt
}

func (o *Order) SetCanceled() error {
	o.Status = OrderStatusCanceled

	return Db.Save(o).Error
}

func (o *Order) SetExpired() {
	o.Status = OrderStatusExpired

	Db.Save(o)
}

func (o *Order) SetSuccess() {
	o.Status = OrderStatusSuccess

	Db.Save(o)
}

func (o *Order) SetFailed() {
	o.Status = OrderStatusFailed

	Db.Save(o)
}

func (o *Order) MarkConfirming(blockNum int64, from, hash string, at time.Time) {
	o.FromAddress = from
	o.ConfirmedAt = &at
	o.RefHash = hash
	o.RefBlockNum = blockNum
	o.Status = OrderStatusConfirming

	Db.Save(o)
}

func (o *Order) SetNotifyState(state int) error {
	o.NotifyNum += 1
	o.NotifyState = state

	return Db.Save(o).Error
}

func (o *Order) GetStatusLabel() string {
	label := "🟢收款成功"
	if o.Status == OrderStatusExpired {
		label = "🔴交易过期"
	}
	if o.Status == OrderStatusWaiting {
		label = "🟡等待支付"
	}
	if o.Status == OrderStatusCanceled {
		label = "⚪️订单取消"
	}

	return label
}

func (o *Order) GetStatusEmoji() string {
	label := "🟢"
	if o.Status == OrderStatusExpired {
		label = "🔴"
	}
	if o.Status == OrderStatusWaiting {
		label = "🟡"
	}
	if o.Status == OrderStatusCanceled {
		label = "⚪️"
	}

	return label
}

func (o *Order) GetTxUrl() string {
	return GetTxUrl(o.TradeType, o.RefHash)
}

func (o *Order) TableName() string {
	return "bep_order"
}

func GetTradeOrder(tradeId string) (Order, bool) {
	var order Order
	res := Db.Where("trade_id = ?", tradeId).Take(&order)

	return order, res.Error == nil
}

func GetOrderByStatus(Status int) []Order {
	orders := make([]Order, 0)

	Db.Where("status = ?", Status).Find(&orders)

	return orders
}

func GetNotifyFailedTradeOrders() ([]Order, error) {
	var orders []Order
	res := Db.Where("status = ?", OrderStatusSuccess).
		Where("notify_num <= ?", GetC(NotifyMaxRetry)).
		Where("notify_state = ?", OrderNotifyStateFail).Find(&orders)

	return orders, res.Error
}

// CalcTradeAmount 计算当前实际可用的交易金额
func CalcTradeAmount(address []string, rate, money decimal.Decimal, t TradeType) (string, string, error) {
	calcMutex.Lock()
	defer calcMutex.Unlock()

	var orders []Order
	lock := make(map[string]bool)
	status := []int{OrderStatusConfirming, OrderStatusWaiting}
	Db.Where("status in (?) and trade_type = ?", status, t).Find(&orders)
	for _, order := range orders {
		lock[order.Address+order.Amount] = true
	}

	atom, precision := GetAtomicity(t)
	if rate.LessThanOrEqual(decimal.Zero) || precision <= 0 {
		return "", "", errors.New(fmt.Sprintf("[%v - %v]原子颗粒度计算异常，联系管理员处理！", atom, precision))
	}

	amount := money.DivRound(rate, precision)
	if amount.LessThan(atom) { // 低于最小原子精度，从最小原子精度开始计算
		amount = atom
	}

	for {
		for _, addr := range address {
			k := addr + amount.String()
			if _, ok := lock[k]; ok {
				continue
			}

			return addr, amount.String(), nil
		}

		// 已经被占用，每次递增一个原子精度
		amount = amount.Add(atom)
	}
}

// CalcTradeExpiredAt 计算订单过期时间 最小180，最大3600，默认1200
func CalcTradeExpiredAt(sec int64) time.Time {
	if sec >= 180 && sec <= 3600 {
		return time.Now().Add(time.Duration(sec) * time.Second)
	}

	return time.Now().Add(time.Duration(cast.ToUint64(GetK(PaymentTimeout))) * time.Second)
}

func GetAtomicity(t TradeType) (decimal.Decimal, int32) {
	confKey, ok := GetTradeAtomKey(t)
	if !ok {
		confKey = "atom_usdt"
	}

	atom, _ := decimal.NewFromString(GetK(confKey))

	return atom, cast.ToInt32(math.Abs(float64(atom.Exponent())))
}
