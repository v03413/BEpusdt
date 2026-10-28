package model

import (
	"time"

	"gorm.io/gorm/clause"
)

// SolanaCursor 记录按收款代币账户检测付款的进度，重启后从这里继续
type SolanaCursor struct {
	Account   string    `gorm:"column:account;type:varchar(64);primaryKey;comment:代币账户" json:"account"`
	Owner     string    `gorm:"column:owner;type:varchar(64);not null;default:'';comment:钱包地址" json:"owner"`
	TradeType string    `gorm:"column:trade_type;type:varchar(20);not null;default:'';comment:交易类型" json:"trade_type"`
	Slot      int64     `gorm:"column:slot;not null;default:0;comment:低水位，slot 小于它的签名都已处理" json:"slot"`
	Signature string    `gorm:"column:signature;type:varchar(128);not null;default:'';comment:最近处理的签名" json:"signature"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;comment:最后更新时间" json:"updated_at"`
}

func (SolanaCursor) TableName() string {

	return "bep_solana_cursor"
}

func GetSolanaCursor(account string) (SolanaCursor, error) {
	var c SolanaCursor
	err := Db.Where("account = ?", account).Limit(1).Find(&c).Error

	return c, err
}

func SaveSolanaCursor(c SolanaCursor) error {
	c.UpdatedAt = time.Now()

	return Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account"}},
		DoUpdates: clause.AssignmentColumns([]string{"owner", "trade_type", "slot", "signature", "updated_at"}),
	}).Create(&c).Error
}
