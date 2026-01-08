package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

func arbitrumInit() {
	ctx := context.Background()
	arb := evm{
		Network: conf.Arbitrum,
		Block: block{
			InitStartOffset: -600,
			ConfirmedOffset: 40,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: arb.blockDispatch})
	Register(Task{Callback: arb.blockRoll, Duration: time.Second * 5})
	Register(Task{Callback: arb.tradeConfirmHandle, Duration: time.Second * 5})
}
