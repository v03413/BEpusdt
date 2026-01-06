package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app"
)

var Version = &cli.Command{
	Name:  "version",
	Usage: "显示版本信息",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("BEpusdt 版本：" + app.Version)

		return nil
	},
}
