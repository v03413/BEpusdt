package command

import "github.com/urfave/cli/v3"

var SQLiteLogFlag = &cli.StringFlag{
	Name:  "sqlite",
	Value: "/var/lib/bepusdt/sqlite.db",
	Usage: "SQLite 数据库文件路径",
}

var LogFlag = &cli.StringFlag{
	Name:  "log",
	Value: "/var/log/bepusdt/",
	Usage: "日志文件保存路径",
}

var ListenFlag = &cli.StringFlag{
	Name:  "listen",
	Value: ":8080",
	Usage: "监听地址，格式为 ip:port，例如 :8080",
}
