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

	dsn := fmt.Sprintf("%s?cache=shared&mode=rwc&_pragma=cache_size(-18888)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path)

	Db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {

		return fmt.Errorf("数据库初始化失败：%w", err)
	}

	{
		sqlDB, err := Db.DB()
		if err != nil {

			return fmt.Errorf("获取数据库连接失败：%w", err)
		}
		sqlDB.SetMaxOpenConns(5)
		sqlDB.SetMaxIdleConns(3)
		sqlDB.SetConnMaxLifetime(0)
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

func Close() {
	if Db == nil {

		return
	}

	sqlDB, err := Db.DB()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("数据库资源句柄获取异常：%s", err.Error()))

		return
	}

	if err := sqlDB.Close(); err != nil {

		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("数据库资源关闭错误：%s", err.Error()))
	}
}
