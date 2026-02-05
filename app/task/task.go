package task

import (
	"context"
	"sync"
	"time"

	"github.com/v03413/bepusdt/app/model"
)

// 区块扫描队列最大长度，避免可能因为 Rpc Rate Limit 问题导致消费队列堆积，进而导致OOM，暂时简单限制队列长度
// 如果直接使用固定长度的 Channel 控制，会导致区块高度同步时也彻底阻塞，无法对外界输出日志导致无法观察，彻底垮掉
// 如果确实是因为 Rate Limit 问题导致的异常，优先考虑的是提升 Rpc 节点的质量和稳定性
const blockQueueLimit = 100

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
	plasmaInit()
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
