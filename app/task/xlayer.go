package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

func xlayerInit() {
	ctx := context.Background()
	xlayer := evm{
		Network:          conf.Xlayer,
		LookbackInterval: time.Millisecond * 300,
		Block: block{
			RollDelayOffset: 3,
			ConfirmedOffset: 12,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: xlayer.blockDispatch})
	Register(Task{Callback: xlayer.syncBlocksForward, Duration: time.Second * 3})
	Register(Task{Callback: xlayer.tradeConfirmHandle, Duration: time.Second * 5})
	Register(Task{Callback: xlayer.lookbackBlocks, Duration: time.Second * 15})
}
