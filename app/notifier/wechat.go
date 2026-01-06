package notifier

import (
	"github.com/v03413/bepusdt/app/model"
)

type Wechat struct{}

func (Wechat) Initialize(params string) error {
	return nil
}

func (Wechat) Success(order model.Order) {}

func (Wechat) NotifyFail(order model.Order, reason string) {}

func (Wechat) NonOrderTransfer(trans model.TronTransfer, wa model.Wallet) {

}

func (Wechat) TronResourceChange(res model.TronResource) {

}

func (Wechat) Welcome() {

}

func (Wechat) Test() error {
	return nil
}
