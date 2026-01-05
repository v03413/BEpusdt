package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/cmd"
	"github.com/v03413/bepusdt/app/conf"
)

func main() {
	c := &cli.Command{
		Name:  "BEpusdt",
		Usage: conf.Desc,
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
