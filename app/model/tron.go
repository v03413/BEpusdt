package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/v03413/tronprotocol/core"
)

type TronTransfer struct {
	Network     string
	TxHash      string
	Amount      decimal.Decimal
	FromAddress string
	RecvAddress string
	Timestamp   time.Time
	TradeType   TradeType
	BlockNum    int
}

type TronResource struct {
	ID           string
	Type         core.Transaction_Contract_ContractType
	Balance      int64
	FromAddress  string
	RecvAddress  string
	Timestamp    time.Time
	ResourceCode core.ResourceCode
}
