package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// 202607081430 - 删除 bep_order.trade_type_reselect 列
func m202607081430DropOrderTradeTypeReselect() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "202607081430_drop_order_trade_type_reselect",
		Migrate: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("bep_order") {
				return nil
			}
			if tx.Migrator().HasColumn("bep_order", "trade_type_reselect") {
				return tx.Exec("ALTER TABLE bep_order DROP COLUMN trade_type_reselect").Error
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("bep_order") {
				return nil
			}
			if !tx.Migrator().HasColumn("bep_order", "trade_type_reselect") {
				return tx.Exec("ALTER TABLE bep_order ADD COLUMN trade_type_reselect INTEGER NOT NULL DEFAULT 0").Error
			}
			return nil
		},
	}
}
