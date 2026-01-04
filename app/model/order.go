package model

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type TradeType string

const (
	OrderNotifyStateSucc = 1 // 回调成功
	OrderNotifyStateFail = 0 // 回调失败

	OrderStatusWaiting    = 1 // 等待支付
	OrderStatusSuccess    = 2 // 交易确认成功
	OrderStatusExpired    = 3 // 订单过期
	OrderStatusCanceled   = 4 // 订单取消
	OrderStatusConfirming = 5 // 等待交易确认
	OrderStatusFailed     = 6 // 交易确认失败

	BscBnb       TradeType = "bsc.bnb"
	EthereumEth  TradeType = "ethereum.eth"
	SolanaSol    TradeType = "solana.sol"
	PolygonMatic TradeType = "polygon.matic"
	AptosApt     TradeType = "aptos.apt"
	ArbitrumEth  TradeType = "arbitrum.eth"
	BaseEth      TradeType = "base.eth"
	XlayerOkb    TradeType = "xlayer.okb"
	TronTrx      TradeType = "tron.trx"

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
)

const (
	OrderApiTypeEpusdt = "epusdt" // epusdt
	OrderApiTypeEpay   = "epay"   // 彩虹易支付
)

var calcMutex sync.Mutex

var explorerUrlMap = map[TradeType]string{
	// Ethereum
	UsdtErc20: "https://etherscan.io/tx/",
	UsdcErc20: "https://etherscan.io/tx/",
	// BSC
	UsdtBep20: "https://bscscan.com/tx/",
	UsdcBep20: "https://bscscan.com/tx/",
	// X Layer
	UsdtXlayer: "https://web3.okx.com/zh-hans/explorer/x-layer/tx/",
	UsdcXlayer: "https://web3.okx.com/zh-hans/explorer/x-layer/tx/",
	// Polygon
	UsdtPolygon: "https://polygonscan.com/tx/",
	UsdcPolygon: "https://polygonscan.com/tx/",
	// Arbitrum
	UsdtArbitrum: "https://arbiscan.io/tx/",
	UsdcArbitrum: "https://arbiscan.io/tx/",
	// Base
	UsdcBase: "https://basescan.org/tx/",
	// Solana
	UsdtSolana: "https://solscan.io/tx/",
	UsdcSolana: "https://solscan.io/tx/",
	// Aptos
	UsdtAptos: "https://explorer.aptoslabs.com/txn/",
	UsdcAptos: "https://explorer.aptoslabs.com/txn/",
	// Tron
	TronTrx:   "https://tronscan.org/#/transaction/",
	UsdtTrc20: "https://tronscan.org/#/transaction/",
	UsdcTrc20: "https://tronscan.org/#/transaction/",
}

var cryptoAtomKeys = map[Crypto]ConfKey{
	USDT: AtomUSDT,
	USDC: AtomUSDC,
	TRX:  AtomTRX,
	BNB:  AtomBNB,
	ETH:  AtomETH,
}

type Order struct {
	Id
	OrderId     string    `Gorm:"column:order_id;type:varchar(128);not null;index;comment:商户ID" json:"order_id"`
	TradeId     string    `Gorm:"column:trade_id;type:varchar(128);not null;uniqueIndex;comment:本地ID" json:"trade_id"`
	TradeType   TradeType `Gorm:"column:trade_type;type:varchar(20);not null;comment:交易类型" json:"trade_type"`
	Fiat        Fiat      `Gorm:"column:fiat;type:varchar(16);not null;index;default:CNY;comment:法币" json:"fiat"`
	Rate        string    `Gorm:"column:rate;type:varchar(10);not null;comment:交易汇率" json:"rate"`
	Amount      string    `Gorm:"column:amount;type:varchar(32);not null;default:0.00;comment:交易数额" json:"amount"`
	Money       string    `Gorm:"column:money;type:varchar(32);not null;default:0.00;comment:订单交易金额" json:"money"`
	Address     string    `Gorm:"column:address;type:varchar(64);not null;comment:收款地址" json:"address"`
	FromAddress string    `Gorm:"type:varchar(34);not null;default:'';comment:支付地址" json:"from_address"`
	Status      int       `Gorm:"type:tinyint(1);not null;default:1;index;comment:交易状态" json:"status"`
	Name        string    `Gorm:"type:varchar(64);not null;default:'';comment:商品名称" json:"name"`
	ApiType     string    `Gorm:"type:varchar(20);not null;default:'epusdt';comment:API类型" json:"api_type"`
	ReturnUrl   string    `Gorm:"type:varchar(255);not null;default:'';comment:同步地址" json:"return_url"`
	NotifyUrl   string    `Gorm:"type:varchar(255);not null;default:'';comment:异步地址" json:"notify_url"`
	NotifyNum   int       `Gorm:"column:notify_num;type:int(11);not null;default:0;comment:回调次数" json:"notify_num"`
	NotifyState int       `Gorm:"column:notify_state;type:tinyint(1);not null;default:0;comment:回调状态 1：成功 0：失败" json:"notify_state"`
	RefHash     string    `Gorm:"column:ref_hash;type:varchar(128);not null;default:'';unique;comment:交易哈希" json:"ref_hash"`
	RefBlockNum int64     `Gorm:"column:ref_block_num;type:bigint(20);not null;default:0;comment:区块索引" json:"ref_block_num"`
	ExpiredAt   time.Time `Gorm:"column:expired_at;type:timestamp;not null;comment:失效时间" json:"expired_at"`
	ConfirmedAt time.Time `Gorm:"type:timestamp;null;comment:交易确认时间"`
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
	o.ConfirmedAt = at
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
	var label = "🟢收款成功"
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
	var label = "🟢"
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

func (o *Order) GetDetailUrl() string {

	return GetDetailUrl(o.TradeType, o.RefHash)
}

func (o *Order) TableName() string {

	return "bep_order"
}

func GetDetailUrl(t TradeType, hash string) string {
	if baseUrl, ok := explorerUrlMap[t]; ok {
		if t == UsdtAptos || t == UsdcAptos {

			return fmt.Sprintf("%s%s?network=mainnet", baseUrl, hash)
		}

		return baseUrl + hash
	}

	return "https://tronscan.org/#/transaction/" + hash
}

func GetTradeOrder(tradeId string) (Order, bool) {
	var order Order
	var res = Db.Where("trade_id = ?", tradeId).Take(&order)

	return order, res.Error == nil
}

func GetOrderByStatus(Status int) []Order {
	var orders = make([]Order, 0)

	Db.Where("status = ?", Status).Find(&orders)

	return orders
}

func GetNotifyFailedTradeOrders() ([]Order, error) {
	var orders []Order
	var res = Db.Where("status = ?", OrderStatusSuccess).
		Where("notify_num <= ?", GetC(NotifyMaxRetry)).
		Where("notify_state = ?", OrderNotifyStateFail).Find(&orders)

	return orders, res.Error
}

// CalcTradeAmount 计算当前实际可用的交易金额
func CalcTradeAmount(address []string, rate, money decimal.Decimal, t TradeType) (string, string) {
	calcMutex.Lock()
	defer calcMutex.Unlock()

	var orders []Order
	var lock = make(map[string]bool)
	var status = []int{OrderStatusConfirming, OrderStatusWaiting}
	Db.Where("status in (?) and trade_type = ?", status, t).Find(&orders)
	for _, order := range orders {
		lock[order.Address+order.Amount] = true
	}

	var atom, precision = getAtomicity(t)
	var amount = money.DivRound(rate, precision)
	if amount.LessThan(atom) { // 低于最小原子精度，从最小原子精度开始计算
		amount = atom
	}

	for {
		for _, addr := range address {
			k := addr + amount.String()
			if _, ok := lock[k]; ok {

				continue
			}

			return addr, amount.String()
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

func getAtomicity(t TradeType) (decimal.Decimal, int32) {
	crypto, ok := TradeTypeTable[t]
	if !ok {
		crypto = USDT
	}

	confKey, ok2 := cryptoAtomKeys[crypto]
	if !ok2 {
		confKey = AtomUSDT
	}

	var atom, _ = decimal.NewFromString(GetK(confKey))

	return atom, cast.ToInt32(math.Abs(float64(atom.Exponent())))
}
