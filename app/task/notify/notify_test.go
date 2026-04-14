package notify

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	applog "github.com/v03413/bepusdt/app/log"
	"github.com/glebarez/sqlite"
	"github.com/v03413/bepusdt/app/model"
	"gorm.io/gorm"
)

func newNotifyTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "notify-test.db")
	db, err := gorm.Open(sqlite.Open(dbPath+"?cache=shared&mode=rwc&_pragma=journal_mode(WAL)&_pragma=busy_timeout(1000)"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	if err := db.AutoMigrate(&model.Order{}); err != nil {
		t.Fatalf("auto migrate order: %v", err)
	}

	return db
}

func initNotifyTestLog(t *testing.T) {
	t.Helper()

	if err := applog.Init(filepath.Join(t.TempDir(), "logs")); err != nil {
		t.Fatalf("init log: %v", err)
	}
	t.Cleanup(func() {
		applog.Close()
	})
}

func newWaitingOrder(notifyURL string) model.Order {
	now := time.Now()
	confirmedAt := now

	return model.Order{
		OrderId:       "merchant-order-1",
		TradeId:       "trade-order-1",
		TradeType:     model.UsdtTrc20,
		Fiat:          "CNY",
		Crypto:        "USDT",
		Rate:          "7.00",
		Amount:        "1.00",
		Money:         "7.00",
		Address:       "TTestAddress1234567890",
		Status:        model.OrderStatusWaiting,
		ApiType:       model.OrderApiTypeEpusdt,
		NotifyUrl:     notifyURL,
		ExpiredAt:     now.Add(10 * time.Minute),
		ConfirmedAt:   &confirmedAt,
		AutoTimeAt:    model.AutoTimeAt{CreatedAt: (*model.Datetime)(&now), UpdatedAt: (*model.Datetime)(&now)},
	}
}

func TestDeliverBepusdtStatusUpdateDoesNotHoldDBWhileHTTPIsPending(t *testing.T) {
	db := newNotifyTestDB(t)
	initNotifyTestLog(t)

	requestStarted := make(chan struct{}, 1)
	releaseResponse := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStarted <- struct{}{}
		<-releaseResponse
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	order := newWaitingOrder(server.URL)
	if err := db.Create(&order).Error; err != nil {
		t.Fatalf("seed order: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- deliverBepusdtStatusUpdate(db, &http.Client{Timeout: 2 * time.Second}, "test-auth-token", order)
	}()

	select {
	case <-requestStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("notification request never reached test server")
	}

	queryDone := make(chan error, 1)
	go func() {
		var count int64
		queryDone <- db.Model(&model.Order{}).Where("status = ?", model.OrderStatusWaiting).Count(&count).Error
	}()

	select {
	case err := <-queryDone:
		if err != nil {
			t.Fatalf("concurrent query failed: %v", err)
		}
	case <-time.After(300 * time.Millisecond):
		close(releaseResponse)
		<-errCh
		t.Fatal("database query was blocked while notification HTTP request was in flight")
	}

	close(releaseResponse)
	if err := <-errCh; err != nil {
		t.Fatalf("deliver notification: %v", err)
	}
}
