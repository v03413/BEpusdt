package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
)

func xlayerInit() {
	ctx := context.Background()
	xlayer := evm{
		Network: conf.Xlayer,
		Block: block{
			InitStartOffset: -600,
			RollDelayOffset: 3,
			ConfirmedOffset: 12,
		},
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: xlayer.blockDispatch})
	Register(Task{Callback: xlayer.blockRoll, Duration: time.Second * 3})
	Register(Task{Callback: xlayer.tradeConfirmHandle, Duration: time.Second * 5})
}
