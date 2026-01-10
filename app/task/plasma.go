package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

func plasmaInit() {
	ctx := context.Background()
	pol := evm{
		Network: conf.Plasma,
		Block: block{
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: pol.blockDispatch})
	Register(Task{Callback: pol.syncBlocksForward, Duration: time.Second * 5})
	Register(Task{Callback: pol.tradeConfirmHandle, Duration: time.Second * 5})
}
