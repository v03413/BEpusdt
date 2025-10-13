package task

import (
	"context"
	"time"

	"github.com/v03413/bepusdt/app/model"
)

type rConf struct {
}

// init 配置刷新到内存，高频读取
func init() {
	var c = rConf{}

	Register(Task{Duration: time.Second * 3, Callback: c.Refresh})
}

func (rConf) Refresh(ctx context.Context) {
	model.RefreshC()
}
