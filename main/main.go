package main

import (
	"context"
	"os"
	"log"

	"github.com/urfave/cli/v3"
	"github.com/joho/godotenv"
	"github.com/v03413/bepusdt/app/cmd"
	"github.com/v03413/bepusdt/app/conf"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, default environment variables will be used.")
	}

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
