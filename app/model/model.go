package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

type Id struct {
	ID int64 `gorm:"column:id;type:INTEGER PRIMARY KEY AUTOINCREMENT;;not null;comment:主键ID" json:"id"`
}

type AutoTimeAt struct {
	CreatedAt *Datetime `gorm:"column:created_at;type:Datetime;not null;comment:记录创建时间;index" json:"created_at"`
	UpdatedAt *Datetime `gorm:"column:updated_at;type:Datetime;not null;comment:最后更新时间" json:"updated_at"`
}

func Init(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {

		return fmt.Errorf("创建数据库目录失败：%w", err)
	}

	Db, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {

		return fmt.Errorf("数据库初始化失败：%w", err)
	}

	if err = AutoMigrate(); err != nil {

		return fmt.Errorf("数据库结构迁移失败：%w", err)
	}

	var count int64
	Db.Model(&Conf{}).Count(&count)
	if count == 0 {
		ConfInit()
	}

	RefreshC()

	return nil
}

func AutoMigrate() error {
	return Db.AutoMigrate(&Wallet{}, &Order{}, &NotifyRecord{}, &Conf{}, &Rate{})
}
