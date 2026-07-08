package cmd

import "github.com/urfave/cli/v3"

var SQLiteFlag = &cli.StringFlag{
	Name:    "sqlite",
	Value:   "/var/lib/bepusdt/sqlite.db",
	Usage:   "SQLite 数据库文件路径",
	Sources: cli.EnvVars("SQLITE"),
}

var PostgresDSNFlag = &cli.StringFlag{
	Name:    "postgres",
	Value:   "",
	Usage:   "PostgreSQL 数据库连接字符串 (DSN)，例如: postgres://user:password@localhost:5432/bepusdt?sslmode=disable&connect_timeout=3",
	Sources: cli.EnvVars("POSTGRESQL_DSN"),
}

var LogFlag = &cli.StringFlag{
	Name:    "log",
	Value:   "/var/log/bepusdt/",
	Usage:   "日志文件保存路径",
	Sources: cli.EnvVars("LOG"),
}

var ListenFlag = &cli.StringFlag{
	Name:    "listen",
	Value:   ":8080",
	Usage:   "监听地址，格式为 ip:port，例如 :8080",
	Sources: cli.EnvVars("LISTEN"),
}
