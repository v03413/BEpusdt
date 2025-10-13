package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
)

func ethInit() {
	ctx := context.Background()
	eth := evm{
		Network: conf.Ethereum,
		Block: block{
			InitStartOffset: -100,
			ConfirmedOffset: 12,
		},
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: eth.blockDispatch})
	Register(Task{Callback: eth.blockRoll, Duration: time.Second * 12})
	Register(Task{Callback: eth.tradeConfirmHandle, Duration: time.Second * 5})
}
