package notifier

import (
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

func TestNonOrderTransferTitleUsesMatchAddress(t *testing.T) {
	trans := model.TronTransfer{
		RecvAddress: "0xf9dff4e813644edc7d9b1b535100b7e92333f743",
		FromAddress: "0x6cba0ff7e76b225ee3b975c3096d828f012c8c72",
	}
	wallet := model.Wallet{
		Address:   "0xF9DFF4E813644eDc7D9b1b535100b7e92333f743",
		MatchAddr: "0xf9dff4e813644edc7d9b1b535100b7e92333f743",
	}

	if title := nonOrderTransferTitle(trans, wallet); title != "收入" {
		t.Fatalf("title = %q, want %q", title, "收入")
	}
}

func TestNonOrderTransferTitleMarksOutgoingTransfer(t *testing.T) {
	trans := model.TronTransfer{
		RecvAddress: "0xrecipient",
		FromAddress: "0xf9dff4e813644edc7d9b1b535100b7e92333f743",
	}
	wallet := model.Wallet{
		Address:   "0xF9DFF4E813644eDc7D9b1b535100b7e92333f743",
		MatchAddr: "0xf9dff4e813644edc7d9b1b535100b7e92333f743",
	}

	if title := nonOrderTransferTitle(trans, wallet); title != "支出" {
		t.Fatalf("title = %q, want %q", title, "支出")
	}
}
