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

type TokenType string

const (
	TokenTypeUSDT TokenType = "USDT"
	TokenTypeUSDC TokenType = "USDC"
	TokenTypeTRX  TokenType = "TRX"
)

// SupportTradeTypes 目前支持的收款交易类型
var SupportTradeTypes = []string{
	TradeTypeTronTrx,
	TradeTypeUsdtTrc20,
	TradeTypeUsdtErc20,
	TradeTypeUsdtBep20,
	TradeTypeUsdtAptos,
	TradeTypeUsdtXlayer,
	TradeTypeUsdtSolana,
	TradeTypeUsdtPolygon,
	TradeTypeUsdtArbitrum,
	TradeTypeUsdcErc20,
	TradeTypeUsdcBep20,
	TradeTypeUsdcXlayer,
	TradeTypeUsdcPolygon,
	TradeTypeUsdcArbitrum,
	TradeTypeUsdcBase,
	TradeTypeUsdcTrc20,
	TradeTypeUsdcSolana,
	TradeTypeUsdcAptos,
}

var TradeTypeTable = map[string]TokenType{
	// USDT
	TradeTypeUsdtTrc20:    TokenTypeUSDT,
	TradeTypeUsdtErc20:    TokenTypeUSDT,
	TradeTypeUsdtBep20:    TokenTypeUSDT,
	TradeTypeUsdtAptos:    TokenTypeUSDT,
	TradeTypeUsdtXlayer:   TokenTypeUSDT,
	TradeTypeUsdtSolana:   TokenTypeUSDT,
	TradeTypeUsdtPolygon:  TokenTypeUSDT,
	TradeTypeUsdtArbitrum: TokenTypeUSDT,

	// USDC
	TradeTypeUsdcErc20:    TokenTypeUSDC,
	TradeTypeUsdcBep20:    TokenTypeUSDC,
	TradeTypeUsdcXlayer:   TokenTypeUSDC,
	TradeTypeUsdcPolygon:  TokenTypeUSDC,
	TradeTypeUsdcArbitrum: TokenTypeUSDC,
	TradeTypeUsdcBase:     TokenTypeUSDC,
	TradeTypeUsdcTrc20:    TokenTypeUSDC,
	TradeTypeUsdcSolana:   TokenTypeUSDC,
	TradeTypeUsdcAptos:    TokenTypeUSDC,

	// TRX
	TradeTypeTronTrx: TokenTypeTRX,
}

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
	if utils.InStrings(wa.TradeType, []string{TradeTypeTronTrx, TradeTypeUsdtTrc20, TradeTypeUsdcTrc20}) {

		return utils.IsValidTronAddress(wa.Address)
	}
	if utils.InStrings(wa.TradeType, []string{TradeTypeUsdtSolana, TradeTypeUsdcSolana}) {

		return utils.IsValidSolanaAddress(wa.Address)
	}
	if utils.InStrings(wa.TradeType, []string{TradeTypeUsdtAptos, TradeTypeUsdcAptos}) {

		return utils.IsValidAptosAddress(wa.Address)
	}

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
	switch wa.TradeType {
	case TradeTypeUsdtPolygon:
		return conf.UsdtPolygon
	case TradeTypeUsdtArbitrum:
		return conf.UsdtArbitrum
	case TradeTypeUsdtErc20:
		return conf.UsdtErc20
	case TradeTypeUsdtBep20:
		return conf.UsdtBep20
	case TradeTypeUsdtXlayer:
		return conf.UsdtXlayer
	case TradeTypeUsdtAptos:
		return conf.UsdtAptos
	case TradeTypeUsdtSolana:
		return conf.UsdtSolana
	case TradeTypeUsdcErc20:
		return conf.UsdcErc20
	case TradeTypeUsdcBep20:
		return conf.UsdcBep20
	case TradeTypeUsdcXlayer:
		return conf.UsdcXlayer
	case TradeTypeUsdcPolygon:
		return conf.UsdcPolygon
	case TradeTypeUsdcArbitrum:
		return conf.UsdcArbitrum
	case TradeTypeUsdcBase:
		return conf.UsdcBase
	case TradeTypeUsdcAptos:
		return conf.UsdcAptos
	case TradeTypeUsdcSolana:
		return conf.UsdcSolana
	default:
		return ""
	}
}

func (wa *Wallet) GetTokenDecimals() int32 {
	switch wa.TradeType {
	case TradeTypeUsdtPolygon:
		return conf.UsdtPolygonDecimals
	case TradeTypeUsdtArbitrum:
		return conf.UsdtArbitrumDecimals
	case TradeTypeUsdtErc20:
		return conf.UsdtEthDecimals
	case TradeTypeUsdtBep20:
		return conf.UsdtBscDecimals
	case TradeTypeUsdtAptos:
		return conf.UsdtAptosDecimals
	case TradeTypeUsdtXlayer:
		return conf.UsdtXlayerDecimals
	case TradeTypeUsdtSolana:
		return conf.UsdtSolanaDecimals
	case TradeTypeUsdcErc20:
		return conf.UsdcEthDecimals
	case TradeTypeUsdcBep20:
		return conf.UsdcBscDecimals
	case TradeTypeUsdcXlayer:
		return conf.UsdcXlayerDecimals
	case TradeTypeUsdcPolygon:
		return conf.UsdcPolygonDecimals
	case TradeTypeUsdcArbitrum:
		return conf.UsdcArbitrumDecimals
	case TradeTypeUsdcBase:
		return conf.UsdcBaseDecimals
	case TradeTypeUsdcSolana:
		return conf.UsdcSolanaDecimals
	case TradeTypeUsdcAptos:
		return conf.UsdcAptosDecimals
	default:
		return -6
	}
}

func GetTokenType(tradeType string) (TokenType, error) {
	if f, ok := TradeTypeTable[tradeType]; ok {
		return f, nil
	}
	return "", fmt.Errorf("unsupported trade type: %s", tradeType)
}

func GetAvailableAddress(tradeType string) []string {
	var rows []Wallet

	Db.Where("trade_type = ? and status = ?", tradeType, WaStatusEnable).Find(&rows)

	var wallets = make([]string, 0)
	for _, w := range rows {
		wallets = append(wallets, w.Address)
	}

	return wallets
}
