package model

import (
	"errors"
	"strings"

	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/xssnick/tonutils-go/address"
)

const (
	WaStatusEnable  uint8 = 1
	WaStatusDisable uint8 = 0
	WaOtherEnable   uint8 = 1
	WaOtherDisable  uint8 = 0
)

type Wallet struct {
	Id
	Name        string `gorm:"column:name;type:varchar(32);not null;default:-';comment:名称" json:"name"`
	Status      uint8  `gorm:"column:status;not null;default:1;comment:地址状态" json:"status"`
	Address     string `gorm:"column:address;type:varchar(128);not null;index;comment:钱包地址" json:"address"`
	MatchAddr   string `gorm:"column:match_addr;type:varchar(128);not null;uniqueIndex:idx_address;comment:匹配地址" json:"match_addr"`
	TradeType   string `gorm:"column:trade_type;type:varchar(20);not null;uniqueIndex:idx_address;comment:交易类型" json:"trade_type"`
	OtherNotify uint8  `gorm:"column:other_notify;not null;default:0;comment:其它通知" json:"other_notify"`
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

func (wa *Wallet) Validate() error {
	tradeType := TradeType(wa.TradeType)

	switch tradeType {
	case TronTrx, UsdtTrc20, UsdcTrc20:
		if !utils.IsValidTronAddress(wa.Address) {
			return errors.New("钱包地址格式不合法，请检查")
		}
	case UsdtSolana, UsdcSolana:
		if !utils.IsValidSolanaAddress(wa.Address) {
			return errors.New("钱包地址格式不合法，请检查")
		}
	case UsdtAptos, UsdcAptos:
		if !utils.IsValidAptosAddress(wa.Address) {
			return errors.New("钱包地址格式不合法，请检查")
		}
	case UsdtTon:
		if !strings.HasPrefix(wa.Address, "UQ") {
			return errors.New("TON 地址必须以 UQ 开头")
		}
		owner, err := address.ParseAddr(wa.Address)
		if err != nil {
			return err
		}
		addr, err := utils.GetJettonWalletAddr(utils.NewTonClient(GetC(RpcGlobalConfigUrlTon)), address.MustParseAddr(conf.UsdtTon), owner)
		if err != nil {
			return err
		}
		wa.MatchAddr = addr.Bounce(false).String()
		return nil
	case TonGram:
		if !strings.HasPrefix(wa.Address, "UQ") {
			return errors.New("TON 地址必须以 UQ 开头")
		}
		owner, err := address.ParseAddr(wa.Address)
		if err != nil {
			return err
		}
		wa.MatchAddr = owner.Bounce(false).String()
		return nil
	default:
		if !utils.IsValidEvmAddress(wa.Address) {
			return errors.New("钱包地址格式不合法，请检查")
		}
	}

	if !AddrCaseSens(tradeType) {
		wa.MatchAddr = strings.ToLower(wa.Address)
	}

	return nil
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

func (wa *Wallet) GetNetwork() Network {
	if c, ok := registry[TradeType(wa.TradeType)]; ok {

		return c.Network
	}

	return ""
}

func (wa *Wallet) GetPaymentAddr() string {
	if wa.TradeType == string(UsdtTon) {

		return wa.Address
	}

	return wa.MatchAddr
}

func (wa *Wallet) GetMatchAddr() string {
	return wa.MatchAddr
}

func GetAvailableWallets(t TradeType) []Wallet {
	var wallets = make([]Wallet, 0)

	Db.Where("trade_type = ? and status = ?", t, WaStatusEnable).Find(&wallets)

	return wallets
}

func NewWallet(address string, tradeType TradeType) (Wallet, error) {
	wa := Wallet{Address: address, TradeType: string(tradeType)}
	if err := wa.Validate(); err != nil {
		return Wallet{}, err
	}

	return wa, nil
}
