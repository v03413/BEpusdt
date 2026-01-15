package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

type Id struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;not null;comment:主键ID" json:"id"`
}

type AutoTimeAt struct {
	CreatedAt *Datetime `gorm:"column:created_at;type:Datetime;not null;comment:记录创建时间;index" json:"created_at"`
	UpdatedAt *Datetime `gorm:"column:updated_at;type:Datetime;not null;comment:最后更新时间" json:"updated_at"`
}

func Init(path string) error {
	return InitSQLite(path)
}

func InitSQLite(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {

		return fmt.Errorf("创建数据库目录失败：%w", err)
	}

	dsn := fmt.Sprintf("%s?cache=shared&mode=rwc"+
		"&_pragma=cache_size(-32000)"+ // 32MB 缓存，平衡内存占用
		"&_pragma=journal_mode(WAL)"+
		"&_pragma=busy_timeout(8000)"+ // 8 秒超时，兼顾慢速磁盘
		"&_pragma=synchronous(NORMAL)"+ // NORMAL 模式，性能与安全平衡
		"&_pragma=wal_autocheckpoint(1500)", // 适中的 checkpoint 频率
		path)
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

func InitMySQL(dsn string) error {
	var err error
	Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {

		return fmt.Errorf("数据库初始化失败：%w", err)
	}

	{
		sqlDB, err := Db.DB()
		if err != nil {

			return fmt.Errorf("获取数据库连接失败：%w", err)
		}
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
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
