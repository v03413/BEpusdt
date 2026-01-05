package model

import (
	"github.com/v03413/bepusdt/app/utils"
)

const (
	WaStatusEnable  uint8 = 1
	WaStatusDisable uint8 = 0
	WaOtherEnable   uint8 = 1
	WaOtherDisable  uint8 = 0
)

const (
	USDT Crypto = "USDT"
	USDC Crypto = "USDC"
	TRX  Crypto = "TRX"
	BNB  Crypto = "BNB"
	ETH  Crypto = "ETH"
)

type Wallet struct {
	Id
	Name        string `gorm:"column:name;type:varchar(32);not null;default:-';comment:名称" json:"name"`
	Status      uint8  `gorm:"column:status;type:tinyint(1);not null;default:1;comment:地址状态" json:"status"`
	Address     string `gorm:"column:address;type:varchar(64);not null;uniqueIndex:idx_address;comment:钱包地址" json:"address"`
	TradeType   string `gorm:"column:trade_type;type:varchar(20);not null;uniqueIndex:idx_address;comment:交易类型" json:"trade_type"`
	OtherNotify uint8  `gorm:"column:other_notify;type:tinyint(1);not null;default:0;comment:其它通知" json:"other_notify"`
	Remark      string `gorm:"column:remark;type:varchar(255);not null;default:'';comment:备注" json:"remark"`
	AutoTimeAt
}

func (wa *Wallet) TableName() string {

	return "bep_wallet"
}

func (wa *Wallet) SetStatus(status uint8) {
	wa.Status = status
	Db.Save(wa)
}

func (wa *Wallet) IsValid() bool {
	tradeType := TradeType(wa.TradeType)

	// Tron 地址验证
	if tradeType == TronTrx || tradeType == UsdtTrc20 || tradeType == UsdcTrc20 {
		return utils.IsValidTronAddress(wa.Address)
	}

	// Solana 地址验证
	if tradeType == UsdtSolana || tradeType == UsdcSolana {
		return utils.IsValidSolanaAddress(wa.Address)
	}

	// Aptos 地址验证
	if tradeType == UsdtAptos || tradeType == UsdcAptos {
		return utils.IsValidAptosAddress(wa.Address)
	}

	// 默认使用 EVM 地址验证（Ethereum, BSC, Polygon, Arbitrum, Base, X Layer）
	return utils.IsValidEvmAddress(wa.Address)
}

func (wa *Wallet) SetOtherNotify(notify uint8) {
	wa.OtherNotify = notify

	Db.Save(wa)
}

func (wa *Wallet) Delete() {
	Db.Delete(wa)
}

func (wa *Wallet) GetTokenContract() string {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {

		return c.Contract
	}

	return ""
}

func (wa *Wallet) GetTokenDecimals() int32 {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {

		return c.Decimal
	}

	return -18
}

func GetAvailableAddress(t TradeType) []string {
	var rows []Wallet
	Db.Where("trade_type = ? and status = ?", t, WaStatusEnable).Find(&rows)

	wallets := make([]string, 0, len(rows))
	for _, w := range rows {
		wallets = append(wallets, w.Address)
	}

	return wallets
}
