package model

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/v03413/bepusdt/app/model/migration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var err error

type Id struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;not null;comment:主键ID" json:"id"`
}

type AutoTimeAt struct {
	CreatedAt *Datetime `gorm:"column:created_at;not null;comment:记录创建时间;index" json:"created_at"`
	UpdatedAt *Datetime `gorm:"column:updated_at;not null;comment:最后更新时间" json:"updated_at"`
}

func Init(db, dsn string) error {
	if dsn != "" {
		return initPostgres(dsn)
	}

	return initSqlite(db)
}

func initSqlite(db string) error {
	if err := os.MkdirAll(filepath.Dir(db), os.ModePerm); err != nil {

		return fmt.Errorf("创建数据库目录失败：%w", err)
	}

	dsn := fmt.Sprintf("%s?cache=shared&mode=rwc"+
		"&_pragma=cache_size(-32000)"+ // 32MB 缓存，平衡内存占用
		"&_pragma=journal_mode(WAL)"+
		"&_pragma=busy_timeout(8000)"+ // 8 秒超时，兼顾慢速磁盘
		"&_pragma=synchronous(NORMAL)"+ // NORMAL 模式，性能与安全平衡
		"&_pragma=wal_autocheckpoint(1500)", // 适中的 checkpoint 频率
		db)
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

	FillDefaultConf()
	RefreshC()

	return nil
}

func initPostgres(dsn string) error {
	// 首次启动可能出现 SLOW SQL 告警，这是由于连接池首次连接预热引起的，后续连接将正常

	var err error
	Db, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {

		return err
	}

	{
		sqlDB, err := Db.DB()
		if err != nil {

			return err
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	}

	if err = AutoMigrate(); err != nil {

		return err
	}

	var count int64
	Db.Model(&Conf{}).Count(&count)
	if count == 0 {
		ConfInit()
	}

	FillDefaultConf()
	RefreshC()

	return nil
}

func AutoMigrate() error {
	return migration.Run(Db, []any{
		&Wallet{},
		&Order{},
		&NotifyRecord{},
		&Conf{},
		&Rate{},
	})
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
