package task

import (
	"context"
	"sync"
	"time"

	"github.com/v03413/bepusdt/app/model"
)

type Task struct {
	Duration time.Duration
	Callback func(ctx context.Context)
}

var (
	tasks []Task
	mu    sync.Mutex
)

func Init() error {
	model.RefreshC()

	bscInit()
	ethInit()
	polygonInit()
	arbitrumInit()
	xlayerInit()
	baseInit()

	return nil
}

func Register(t Task) {
	mu.Lock()
	defer mu.Unlock()

	if t.Callback == nil {

		panic("Task Callback cannot be nil")
	}

	tasks = append(tasks, t)
}

func Start(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	for _, t := range tasks {
		go func(t Task) {
			if t.Duration <= 0 {
				t.Callback(ctx)

				return
			}

			t.Callback(ctx)

			ticker := time.NewTicker(t.Duration)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					t.Callback(ctx)
				}
			}
		}(t)
	}
}
