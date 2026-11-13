package task

import (
	"testing"

	"github.com/v03413/bepusdt/app/model"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	tgo "github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

var (
	tonTestSrc = address.NewAddress(0, 0, make([]byte, 32))
	tonTestDst = address.NewAddress(0, 0, append(make([]byte, 31), 1))
)

func tonTestTx(body *cell.Cell, nano uint64, dsc any) *tlb.Transaction {
	tx := &tlb.Transaction{AccountAddr: tonTestDst.Data(), Now: 1700000000, Hash: []byte{0xab}, Description: dsc}
	tx.IO.In = &tlb.Message{MsgType: tlb.MsgTypeInternal, Msg: &tlb.InternalMessage{
		SrcAddr: tonTestSrc,
		DstAddr: tonTestDst,
		Amount:  tlb.FromNanoTONU(nano),
		Body:    body,
	}}

	return tx
}

func tonTestOrdinary(success bool) tlb.TransactionDescriptionOrdinary {
	return tlb.TransactionDescriptionOrdinary{
		ComputePhase: tlb.ComputePhase{Phase: tlb.ComputePhaseVM{Success: success}},
		ActionPhase:  &tlb.ActionPhase{Success: success},
	}
}

func tonTestJettonBody(op uint32, amount uint64) *cell.Cell {
	return cell.BeginCell().
		MustStoreUInt(uint64(op), 32).
		MustStoreUInt(1, 64).
		MustStoreCoins(amount).
		MustStoreAddr(tonTestSrc).
		EndCell()
}

func TestParseTonTransfer(t *testing.T) {
	empty := cell.BeginCell().EndCell()
	skipped := tlb.TransactionDescriptionOrdinary{
		ComputePhase: tlb.ComputePhase{Phase: tlb.ComputePhaseSkipped{}},
		Aborted:      true,
	}
	returned := skipped
	returned.BouncePhase = &tlb.BouncePhase{Phase: tlb.BouncePhaseOk{}}
	noFunds := skipped
	noFunds.BouncePhase = &tlb.BouncePhase{Phase: tlb.BouncePhaseNoFunds{}}

	bouncedIn := tonTestTx(empty, 1e9, tonTestOrdinary(true))
	bouncedIn.IO.In.Msg.(*tlb.InternalMessage).Bounced = true
	external := tonTestTx(empty, 1e9, tonTestOrdinary(true))
	external.IO.In.Msg = &tlb.ExternalMessage{}
	noIn := tonTestTx(empty, 1e9, tonTestOrdinary(true))
	noIn.IO.In = nil

	cases := []struct {
		name string
		tx   *tlb.Transaction
		ok   bool
	}{
		{"plain", tonTestTx(empty, 1e9, tonTestOrdinary(true)), true},
		{"uninit account", tonTestTx(empty, 1e9, skipped), true},
		{"uninit account, bounce unpaid", tonTestTx(empty, 1e9, noFunds), true},
		{"returned to sender", tonTestTx(empty, 1e9, returned), false},
		{"zero amount", tonTestTx(empty, 0, tonTestOrdinary(true)), false},
		{"with body", tonTestTx(tonTestJettonBody(tonOpInternalTransfer, 1), 1e9, tonTestOrdinary(true)), false},
		{"bounced inbound", bouncedIn, false},
		{"external inbound", external, false},
		{"no inbound", noIn, false},
		{"tick tock", tonTestTx(empty, 1e9, tlb.TransactionDescriptionTickTock{}), false},
	}

	for _, c := range cases {
		got, ok := (&ton{}).parseTonTransfer(c.tx, 100)
		if ok != c.ok {
			t.Errorf("%s: ok = %v, want %v", c.name, ok, c.ok)
		}
		if !ok {
			continue
		}
		if got.TradeType != model.TonGram || got.Amount.String() != "1" || got.BlockNum != 100 {
			t.Errorf("%s: unexpected transfer %+v", c.name, got)
		}
		if got.RecvAddress != tonTestDst.Bounce(false).String() || got.FromAddress != tonTestSrc.Bounce(false).String() {
			t.Errorf("%s: unexpected address %s <- %s", c.name, got.RecvAddress, got.FromAddress)
		}
	}
}

func TestParseTonInternalTransfer(t *testing.T) {
	shard := &tgo.BlockIDExt{Workchain: 0}
	cases := []struct {
		name string
		tx   *tlb.Transaction
		ok   bool
	}{
		{"jetton received", tonTestTx(tonTestJettonBody(tonOpInternalTransfer, 2_500_000), 1, tonTestOrdinary(true)), true},
		{"other op", tonTestTx(tonTestJettonBody(0x7362d09c, 2_500_000), 1, tonTestOrdinary(true)), false},
		{"execution failed", tonTestTx(tonTestJettonBody(tonOpInternalTransfer, 2_500_000), 1, tonTestOrdinary(false)), false},
		{"zero amount", tonTestTx(tonTestJettonBody(tonOpInternalTransfer, 0), 1, tonTestOrdinary(true)), false},
		{"empty body", tonTestTx(cell.BeginCell().EndCell(), 1, tonTestOrdinary(true)), false},
	}

	for _, c := range cases {
		got, ok := (&ton{}).parseInternalTransfer(shard, c.tx, 100)
		if ok != c.ok {
			t.Errorf("%s: ok = %v, want %v", c.name, ok, c.ok)
		}
		if !ok {
			continue
		}
		if got.TradeType != model.UsdtTon || got.Amount.String() != "2.5" || got.RecvAddress != tonTestDst.Bounce(false).String() {
			t.Errorf("%s: unexpected transfer %+v", c.name, got)
		}
	}
}
