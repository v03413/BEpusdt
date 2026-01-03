package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/command"
)

func main() {
	cmd := &cli.Command{
		Name:  "BEpusdt",
		Usage: "一款更好用的个人加密货币收款网关",
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
