package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/cmd"
)

func main() {
	c := &cli.Command{
		Name:  "BEpusdt",
		Usage: "一款更好用的个人加密货币收款网关",
		Commands: []*cli.Command{
			cmd.Start,
			cmd.Version,
			cmd.Reset,
		},
	}
	if err := c.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
