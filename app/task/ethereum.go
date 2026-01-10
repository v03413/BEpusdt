package task

import (
	"context"
	"time"

	"github.com/smallnest/chanx"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

func ethInit() {
	ctx := context.Background()
	eth := evm{
		Network: conf.Ethereum,
		Block: block{
			ConfirmedOffset: 12,
		},
		Native: evmNative{
			Parse:     true,
			TradeType: model.EthereumEth,
			Decimal:   conf.EthereumEthDecimals,
		},
		Client:         utils.NewHttpClient(),
		blockScanQueue: chanx.NewUnboundedChan[evmBlock](ctx, 30),
	}

	Register(Task{Callback: eth.blockDispatch})
	Register(Task{Callback: eth.syncBlocksForward, Duration: time.Second * 12})
	Register(Task{Callback: eth.tradeConfirmHandle, Duration: time.Second * 5})
}
