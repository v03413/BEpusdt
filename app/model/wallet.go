package model

import (
	"fmt"

	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

const (
	WaStatusEnable  uint8 = 1
	WaStatusDisable uint8 = 0
	WaOtherEnable   uint8 = 1
	WaOtherDisable  uint8 = 0
)

type Crypto string

const (
	USDT Crypto = "USDT"
	USDC Crypto = "USDC"
	TRX  Crypto = "TRX"
	BNB  Crypto = "BNB"
	ETH  Crypto = "ETH"
)

// SupportTradeTypes 目前支持的收款交易类型
var SupportTradeTypes = map[TradeType]struct{}{
	EthereumEth:  {},
	BscBnb:       {},
	TronTrx:      {},
	UsdtTrc20:    {},
	UsdtErc20:    {},
	UsdtBep20:    {},
	UsdtAptos:    {},
	UsdtXlayer:   {},
	UsdtSolana:   {},
	UsdtPolygon:  {},
	UsdtArbitrum: {},
	UsdcErc20:    {},
	UsdcBep20:    {},
	UsdcXlayer:   {},
	UsdcPolygon:  {},
	UsdcArbitrum: {},
	UsdcBase:     {},
	UsdcTrc20:    {},
	UsdcSolana:   {},
	UsdcAptos:    {},
}

var TradeTypeTable = map[TradeType]Crypto{
	// USDT
	UsdtTrc20:    USDT,
	UsdtErc20:    USDT,
	UsdtBep20:    USDT,
	UsdtAptos:    USDT,
	UsdtXlayer:   USDT,
	UsdtSolana:   USDT,
	UsdtPolygon:  USDT,
	UsdtArbitrum: USDT,

	// USDC
	UsdcErc20:    USDC,
	UsdcBep20:    USDC,
	UsdcXlayer:   USDC,
	UsdcPolygon:  USDC,
	UsdcArbitrum: USDC,
	UsdcBase:     USDC,
	UsdcTrc20:    USDC,
	UsdcSolana:   USDC,
	UsdcAptos:    USDC,

	// 原生代币
	TronTrx:     TRX,
	EthereumEth: ETH,
	BscBnb:      BNB,
}

// tokenContractMap Crypto 合约地址映射表
var tokenContractMap = map[TradeType]string{
	UsdtPolygon:  conf.UsdtPolygon,
	UsdtArbitrum: conf.UsdtArbitrum,
	UsdtErc20:    conf.UsdtErc20,
	UsdtBep20:    conf.UsdtBep20,
	UsdtXlayer:   conf.UsdtXlayer,
	UsdtAptos:    conf.UsdtAptos,
	UsdtSolana:   conf.UsdtSolana,
	UsdcErc20:    conf.UsdcErc20,
	UsdcBep20:    conf.UsdcBep20,
	UsdcXlayer:   conf.UsdcXlayer,
	UsdcPolygon:  conf.UsdcPolygon,
	UsdcArbitrum: conf.UsdcArbitrum,
	UsdcBase:     conf.UsdcBase,
	UsdcAptos:    conf.UsdcAptos,
	UsdcSolana:   conf.UsdcSolana,
}

// tokenDecimalsMap Crypto 精度映射表
var tokenDecimalsMap = map[TradeType]int32{
	UsdtPolygon:  conf.UsdtPolygonDecimals,
	UsdtArbitrum: conf.UsdtArbitrumDecimals,
	UsdtErc20:    conf.UsdtEthDecimals,
	UsdtBep20:    conf.UsdtBscDecimals,
	UsdtAptos:    conf.UsdtAptosDecimals,
	UsdtXlayer:   conf.UsdtXlayerDecimals,
	UsdtSolana:   conf.UsdtSolanaDecimals,
	UsdtTrc20:    conf.UsdtTronDecimals,
	UsdcErc20:    conf.UsdcEthDecimals,
	UsdcBep20:    conf.UsdcBscDecimals,
	UsdcXlayer:   conf.UsdcXlayerDecimals,
	UsdcPolygon:  conf.UsdcPolygonDecimals,
	UsdcArbitrum: conf.UsdcArbitrumDecimals,
	UsdcBase:     conf.UsdcBaseDecimals,
	UsdcSolana:   conf.UsdcSolanaDecimals,
	UsdcAptos:    conf.UsdcAptosDecimals,
	UsdcTrc20:    conf.UsdcTronDecimals,
}

type Wallet struct {
	Id
	Name        string `Gorm:"column:name;type:varchar(32);not null;default:-';comment:名称" json:"name"`
	Status      uint8  `Gorm:"column:status;type:tinyint(1);not null;default:1;comment:地址状态" json:"status"`
	Address     string `Gorm:"column:address;type:varchar(64);not null;uniqueIndex:idx_address;comment:钱包地址" json:"address"`
	TradeType   string `Gorm:"column:trade_type;type:varchar(20);not null;uniqueIndex:idx_address;comment:交易类型" json:"trade_type"`
	OtherNotify uint8  `Gorm:"column:other_notify;type:tinyint(1);not null;default:0;comment:其它通知" json:"other_notify"`
	Remark      string `Gorm:"column:remark;type:varchar(255);not null;default:'';comment:备注" json:"remark"`
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
	if contract, ok := tokenContractMap[TradeType(wa.TradeType)]; ok {

		return contract
	}

	return ""
}

func (wa *Wallet) GetTokenDecimals() int32 {
	if decimals, ok := tokenDecimalsMap[TradeType(wa.TradeType)]; ok {

		return decimals
	}

	return -18
}

func GetCrypto(t TradeType) (Crypto, error) {
	if f, ok := TradeTypeTable[t]; ok {
		return f, nil
	}
	return "", fmt.Errorf("unsupported trade type: %s", t)
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
