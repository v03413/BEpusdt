package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const TableName = "bep_migration"

var migrations = []*gormigrate.Migration{
	m202607081430DropOrderTradeTypeReselect(),
}

func Run(db *gorm.DB, initModels []any) error {
	if err := db.AutoMigrate(initModels...); err != nil {
		return err
	}

	options := &gormigrate.Options{TableName: TableName}

	// 旧版升级/全新安装，构建迁移表
	if !db.Migrator().HasTable(TableName) {
		if err := db.Exec("CREATE TABLE " + TableName + " (id VARCHAR(255) PRIMARY KEY)").Error; err != nil {
			return err
		}
	}

	m := gormigrate.New(db, options, migrations)

	return m.Migrate()
}
