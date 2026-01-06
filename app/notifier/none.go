package notifier

import "github.com/v03413/bepusdt/app/model"

type None struct{}

func (None) Initialize(params string) error {
	return nil
}

func (None) Success(order model.Order) {}

func (None) NotifyFail(order model.Order, reason string) {}

func (None) NonOrderTransfer(trans model.TronTransfer, wa model.Wallet) {

}

func (None) TronResourceChange(res model.TronResource) {

}

func (None) Welcome() {

}

func (None) Test() error {
	return nil
}
