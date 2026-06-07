package utils

import (
	"context"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	tgo "github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
)

func NewTonClient(configUrl string) tgo.APIClientWrapped {
	pool := liteclient.NewConnectionPool()
	if err := pool.AddConnectionsFromConfigUrl(context.Background(), configUrl); err != nil {
		panic(err)
	}

	return tgo.NewAPIClient(pool).WithRetryTimeout(3, time.Second*5)
}

func GetJettonWalletAddr(api tgo.APIClientWrapped, master *address.Address, owner *address.Address) (*address.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client := jetton.NewJettonMasterClient(api, master)
	wallet, err := client.GetJettonWallet(ctx, owner)
	if err != nil {
		return nil, err
	}

	return wallet.Address(), err
}
