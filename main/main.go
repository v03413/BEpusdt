package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/command"
	"github.com/v03413/bepusdt/app/model"
)

func main() {
	cmd := &cli.Command{
		Name:  "BEpusdt",
		Usage: "一款更好用的个人加密货币收款网关",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "sqlite",
				Value: "/var/lib/bepusdt/sqlite.db",
				Usage: "SQLite 数据库文件路径",
			},
			&cli.StringFlag{
				Name:  "log",
				Value: "/var/log/bepusdt.log",
				Usage: "日志文件保存路径",
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			err := model.Init(c.String("sqlite"))
			if err != nil {

				return ctx, fmt.Errorf("数据库初始化失败 %w", err)
			}

			return ctx, nil
		},
		Commands: []*cli.Command{
			command.Start,
			command.Version,
			command.Reset,
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
