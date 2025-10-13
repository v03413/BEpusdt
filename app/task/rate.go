package task

import (
	"context"
	"time"

	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/model"
)

func init() {
	var rate = Rate{}

	Register(Task{Duration: time.Second * 5, Callback: rate.Sync})
}

type Rate struct {
}

func (Rate) Sync(ctx context.Context) {
	var lastAt model.Rate
	model.Db.Model(&lastAt).Order("id desc").Find(&lastAt)
	var interval = cast.ToInt64(model.GetC(model.RateSyncInterval))
	if lastAt.ID != 0 && time.Now().Unix()-lastAt.CreatedAt.Time().Unix() < interval {

		return
	}

	model.CoingeckoRate()
}
