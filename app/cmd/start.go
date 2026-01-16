package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"context"

	"github.com/v03413/bepusdt/app"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/notifier"
	"github.com/v03413/bepusdt/app/router"
	"github.com/v03413/bepusdt/app/task"
)

import (
	"github.com/urfave/cli/v3"
)

var Start = &cli.Command{
	Name:  "start",
	Usage: "启动收款网关",
	Flags: []cli.Flag{SQLiteFlag, MySQLDSNFlag, LogFlag, ListenFlag},
	Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
		dsn := c.String("mysql")
		sqlite := c.String("sqlite")
		if err := model.Init(sqlite, dsn); err != nil {
			return ctx, fmt.Errorf("数据库初始化失败 %w", err)
		}

		if err := log.Init(c.String("log")); err != nil {
			return ctx, fmt.Errorf("日志初始化失败 %w", err)
		}

		return ctx, task.Init()
	},
	After: func(ctx context.Context, c *cli.Command) error {
		log.Close()
		model.Close()

		return nil
	},
	Action: start,
}

func start(ctx context.Context, cmd *cli.Command) error {
	// 开始任务调度
	task.Start(ctx)

	// 启动 Web 服务器
	var listen = cmd.String("listen")
	var srv = &http.Server{Addr: listen, Handler: router.Handler()}

	log.Info("web server Start listen", listen)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("web server error", err)
		}
	}()

	// 关闭 Web 服务器
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdown); err != nil {
			log.Error("Web shutdown Error", err)
			return
		}

		log.Info("web shutdown success.")
	}()

	notifier.Welcome()

	fmt.Println(fmt.Sprintf("BEpusdt 启动成功(%s)，当前版本：%s", listen, app.Version))

	// 等待中断信号
	var signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	runtime.GC()

	return nil
}
