package task

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
)

const (
	testSolanaOwner       = "DAzEQJ8TzdmrAgphrcGGZie4fwiXmRYXRCKX4wQh2oLf"
	testSolanaTokenWallet = "23PXKLkUNQ85LScKDZVWhHLzFppYiSVky2Z7iqkk4JFu"
	testSolanaSender      = "5tzFkiKscXHK5ZXCGbXZxdw7gTjjD1mBwuoFbhUvuAi9"
	testSolanaTxHash      = "3uhzPPDFKz8JTWJTPN2D7VbDTL94KXw2WyqEUvpLvbcsKVymDKETfMgwUxGsWM2tBDvteTHSRNZPU3Cx8u2hyMbe"
)

func TestSolanaReconcileWaitingOrdersMarksMatchingOrderConfirming(t *testing.T) {
	initSolanaReconcileTestLog(t)

	dbPath := filepath.Join(t.TempDir(), "solana-reconcile.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode rpc request: %v", err)
		}

		method, _ := req["method"].(string)
		switch method {
		case "getTokenAccountsByOwner":
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"value":[{"pubkey":"` + testSolanaTokenWallet + `"}]}}`))
		case "getSignaturesForAddress":
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":[{"signature":"` + testSolanaTxHash + `","slot":424729879,"err":null,"blockTime":1780769985}]}`))
		case "getTransaction":
			_, _ = w.Write([]byte(`{
				"jsonrpc":"2.0",
				"id":1,
				"result":{
					"slot":424729879,
					"blockTime":1780769985,
					"meta":{
						"err":null,
						"preTokenBalances":[
							{"accountIndex":2,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaSender + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
							{"accountIndex":3,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaOwner + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"}
						],
						"postTokenBalances":[
							{"accountIndex":2,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaSender + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
							{"accountIndex":3,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaOwner + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"}
						],
						"innerInstructions":[]
					},
					"transaction":{
						"message":{
							"accountKeys":[
								{"pubkey":"` + testSolanaSender + `"},
								{"pubkey":"13gxbc8s6rPLDXPMiTnbnKD6uVdddzKpUCBZA5iCmZkd"},
								{"pubkey":"7KJjY7rArbydeLBF7gQ5LdqXRKRYyPArT99NEctsHsgU"},
								{"pubkey":"` + testSolanaTokenWallet + `"}
							],
							"instructions":[
								{
									"program":"spl-token",
									"programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
									"parsed":{
										"type":"transferChecked",
										"info":{
											"authority":"` + testSolanaSender + `",
											"source":"7KJjY7rArbydeLBF7gQ5LdqXRKRYyPArT99NEctsHsgU",
											"destination":"` + testSolanaTokenWallet + `",
											"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
											"tokenAmount":{"amount":"13190000","decimals":6,"uiAmount":13.19,"uiAmountString":"13.19"}
										}
									}
								}
							]
						}
					}
				}
			}`))
		default:
			t.Fatalf("unexpected rpc method: %s", method)
		}
	}))
	defer server.Close()

	model.SetK(model.RpcEndpointSolana, server.URL)

	createdAt := model.Datetime(time.Date(2026, 6, 7, 2, 13, 6, 0, time.Local))
	expiredAt := time.Date(2026, 6, 7, 2, 43, 6, 0, time.Local)
	zero := time.Unix(0, 0)
	order := model.Order{
		OrderId:     "sub2_20260607ef8AeiFK",
		TradeId:     "csxMKKMccWtlbGWg2Y",
		TradeType:   model.UsdcSolana,
		Crypto:      model.USDC,
		Amount:      "13.19",
		Money:       "100",
		Address:     testSolanaOwner,
		Status:      model.OrderStatusWaiting,
		NotifyUrl:   "https://example.com/webhook",
		ConfirmedAt: &zero,
		ExpiredAt:   expiredAt,
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	s := newSolana()
	s.client = server.Client()
	s.reconcileWaitingOrders(context.Background())

	var refreshed model.Order
	if err := model.Db.First(&refreshed, order.ID).Error; err != nil {
		t.Fatalf("reload order: %v", err)
	}

	if refreshed.Status != model.OrderStatusConfirming {
		t.Fatalf("expected status %d, got %d", model.OrderStatusConfirming, refreshed.Status)
	}
	if refreshed.RefHash != testSolanaTxHash {
		t.Fatalf("expected ref hash %s, got %s", testSolanaTxHash, refreshed.RefHash)
	}
	if refreshed.FromAddress != testSolanaSender {
		t.Fatalf("expected from address %s, got %s", testSolanaSender, refreshed.FromAddress)
	}
	if refreshed.ConfirmedAt == nil {
		t.Fatal("expected confirmed_at to be set")
	}
	if got := refreshed.ConfirmedAt.Unix(); got != 1780769985 {
		t.Fatalf("expected confirmed_at unix %d, got %d", 1780769985, got)
	}
}

func initSolanaReconcileTestLog(t *testing.T) {
	t.Helper()

	logger := logrus.New()
	logger.SetOutput(io.Discard)
	log.Task = logger
}

func TestSolanaParsedTransactionProducesMatchingTransfer(t *testing.T) {
	tx := gjson.Parse(`{
		"blockTime":1780769985,
		"meta":{
			"preTokenBalances":[
				{"accountIndex":2,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaSender + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
				{"accountIndex":3,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaOwner + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"}
			],
			"postTokenBalances":[
				{"accountIndex":2,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaSender + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"},
				{"accountIndex":3,"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","owner":"` + testSolanaOwner + `","programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"}
			],
			"innerInstructions":[]
		},
		"transaction":{
			"message":{
				"accountKeys":[
					{"pubkey":"` + testSolanaSender + `"},
					{"pubkey":"13gxbc8s6rPLDXPMiTnbnKD6uVdddzKpUCBZA5iCmZkd"},
					{"pubkey":"7KJjY7rArbydeLBF7gQ5LdqXRKRYyPArT99NEctsHsgU"},
					{"pubkey":"` + testSolanaTokenWallet + `"}
				],
				"instructions":[
					{
						"program":"spl-token",
						"programId":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
						"parsed":{
							"type":"transferChecked",
							"info":{
								"authority":"` + testSolanaSender + `",
								"source":"7KJjY7rArbydeLBF7gQ5LdqXRKRYyPArT99NEctsHsgU",
								"destination":"` + testSolanaTokenWallet + `",
								"mint":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
								"tokenAmount":{"amount":"13190000","decimals":6,"uiAmount":13.19,"uiAmountString":"13.19"}
							}
						}
					}
				]
			}
		}
	}`)

	transfers := parseSolanaParsedTransfers(tx)
	if len(transfers) != 1 {
		t.Fatalf("expected 1 transfer, got %d", len(transfers))
	}

	got := transfers[0]
	if got.TradeType != model.UsdcSolana {
		t.Fatalf("expected trade type %s, got %s", model.UsdcSolana, got.TradeType)
	}
	if got.FromAddress != testSolanaSender {
		t.Fatalf("expected from %s, got %s", testSolanaSender, got.FromAddress)
	}
	if got.RecvAddress != testSolanaOwner {
		t.Fatalf("expected recv %s, got %s", testSolanaOwner, got.RecvAddress)
	}
	if got.Amount.String() != "13.19" {
		t.Fatalf("expected amount 13.19, got %s", got.Amount.String())
	}
	if got.Timestamp.Unix() != 1780769985 {
		t.Fatalf("expected timestamp 1780769985, got %d", got.Timestamp.Unix())
	}
}

func TestGetConfirmingOrdersKeepsOnTimePaymentAfterOrderExpires(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "confirming-after-expire.db")
	if err := model.Init(dbPath, "", ""); err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(model.Close)

	createdAt := model.Datetime(time.Now().Add(-2 * time.Hour))
	confirmedAt := time.Now().Add(-90 * time.Minute)
	expiredAt := time.Now().Add(-1 * time.Hour)
	order := model.Order{
		OrderId:       "late-finalized",
		TradeId:       "late-finalized-trade",
		TradeType:     model.UsdcSolana,
		Crypto:        model.USDC,
		Amount:        "13.19",
		Money:         "100",
		Address:       testSolanaOwner,
		Status:        model.OrderStatusConfirming,
		RefHash:       testSolanaTxHash,
		RefBlockNum:   424729879,
		NotifyUrl:     "https://example.com/webhook",
		ConfirmedAt:   &confirmedAt,
		ExpiredAt:     expiredAt,
		AddressLocked: false,
		AutoTimeAt: model.AutoTimeAt{
			CreatedAt: &createdAt,
			UpdatedAt: &createdAt,
		},
	}
	if err := model.Db.Create(&order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}

	orders := getConfirmingOrders([]model.TradeType{model.UsdcSolana})
	if len(orders) != 1 {
		t.Fatalf("expected confirming order to remain eligible, got %d", len(orders))
	}

	var refreshed model.Order
	if err := model.Db.First(&refreshed, order.ID).Error; err != nil {
		t.Fatalf("reload order: %v", err)
	}
	if refreshed.Status != model.OrderStatusConfirming {
		t.Fatalf("expected status %d, got %d", model.OrderStatusConfirming, refreshed.Status)
	}
}
