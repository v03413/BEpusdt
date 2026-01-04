package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app"
	epay2 "github.com/v03413/bepusdt/app/handler/epay"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/notifier"
	"github.com/v03413/bepusdt/app/utils"

	"github.com/v03413/go-cache"
)

type EpNotify struct {
	TradeId            string  `json:"trade_id"`             //  本地订单号
	OrderId            string  `json:"order_id"`             //  客户交易id
	Amount             float64 `json:"amount"`               //  订单金额 CNY
	ActualAmount       string  `json:"actual_amount"`        //  USDT 交易数额
	Token              string  `json:"token"`                //  收款钱包地址
	BlockTransactionId string  `json:"block_transaction_id"` // 区块id
	Signature          string  `json:"signature"`            // 签名
	Status             int     `json:"status"`               //  1：等待支付，2：支付成功，3：订单超时
}

func Handle(order model.Order) {
	if order.Status != model.OrderStatusSuccess {

		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if order.ApiType == model.OrderApiTypeEpay {
		epay(ctx, order)

		return
	}

	epusdt(ctx, order)
}

func epay(ctx context.Context, order model.Order) {
	var client = http.Client{Timeout: time.Second * 5}
	var notifyUrl = fmt.Sprintf("%s?%s", order.NotifyUrl, epay2.BuildNotifyParams(order))

	postReq, err2 := http.NewRequestWithContext(ctx, "GET", notifyUrl, nil)
	if err2 != nil {
		log.Error("Notify NewRequest Error：", err2)

		return
	}

	postReq.Header.Set("Powered-By", "https://github.com/v03413/bepusdt")
	resp, err := client.Do(postReq)
	if err != nil {
		log.Error("Notify Handle Error：", err)

		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		markNotifyFail(order, fmt.Sprintf("resp.StatusCode != 200"))

		return
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		markNotifyFail(order, fmt.Sprintf("io.ReadAll(resp.Body) Error: %v", err))

		return
	}

	var bodyStr = strings.ToLower(strings.TrimSpace(string(all)))

	// 判断是否包含 success 或 ok
	if !strings.Contains(bodyStr, "success") && !strings.Contains(bodyStr, "ok") {
		markNotifyFail(order, "商户系统必须响应 success 或 ok 才会认定回调成功")

		return
	}

	if err = order.SetNotifyState(model.OrderNotifyStateSucc); err != nil {
		log.Error("订单标记通知成功错误：", err, order.OrderId)

		return
	}

	log.Info("订单通知成功：", order.OrderId)
}

func epusdt(ctx context.Context, order model.Order) {
	var data = make(map[string]interface{})
	var body = EpNotify{
		TradeId:            order.TradeId,
		OrderId:            order.OrderId,
		Amount:             cast.ToFloat64(order.Money),
		ActualAmount:       order.Amount,
		Token:              order.Address,
		BlockTransactionId: order.RefHash,
		Status:             order.Status,
	}
	var jsonBody, err = json.Marshal(body)
	if err != nil {
		log.Error("Notify Json Marshal Error：", err)

		return
	}

	if err = json.Unmarshal(jsonBody, &data); err != nil {
		log.Error("Notify JSON Unmarshal Error：", err)

		return
	}

	// 签名
	body.Signature = utils.EpusdtSign(data, model.AuthToken())

	// 再次序列化
	jsonBody, err = json.Marshal(body)
	var client = http.Client{Timeout: time.Second * 5}
	var postReq, err2 = http.NewRequestWithContext(ctx, "POST", order.NotifyUrl, strings.NewReader(string(jsonBody)))
	if err2 != nil {
		markNotifyFail(order, err.Error())

		return
	}

	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Powered-By", "https://github.com/v03413/bepusdt")
	postReq.Header.Set("User-Agent", "BEpusdt/"+app.Version)
	resp, err := client.Do(postReq)
	if err != nil {
		markNotifyFail(order, err.Error())

		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		markNotifyFail(order, fmt.Sprintf("商户系统返回状态码错误：%d（必须是200）", resp.StatusCode))

		return
	}

	if err = order.SetNotifyState(model.OrderNotifyStateSucc); err != nil {
		log.Error("订单标记通知成功错误：", err, order.OrderId)

		return
	}

	log.Info("订单通知成功：", order.OrderId)
}

func Bepusdt(order model.Order) {
	if order.ApiType != model.OrderApiTypeEpusdt {

		return
	}

	var todo = func() error {
		var o model.Order
		var db = model.Db.Begin()
		if err := db.Where("trade_id = ? and status = ?", order.TradeId, order.Status).First(&o).Error; err != nil {
			db.Rollback()

			return err
		}

		var key = fmt.Sprintf("bepusdt_notify_%d_%s", o.Status, o.TradeId)
		if _, ok := cache.Get(key); ok {
			db.Rollback()

			return nil
		}

		cache.Set(key, true, time.Minute)

		var data = make(map[string]interface{})
		var body = EpNotify{
			TradeId:            o.TradeId,
			OrderId:            o.OrderId,
			Amount:             cast.ToFloat64(o.Money),
			ActualAmount:       o.Amount,
			Token:              o.Address,
			BlockTransactionId: o.RefHash,
			Status:             o.Status,
		}
		var jsonBody, err = json.Marshal(body)
		if err != nil {
			db.Rollback()

			return err
		}

		if err = json.Unmarshal(jsonBody, &data); err != nil {
			db.Rollback()

			return err
		}

		// 签名
		body.Signature = utils.EpusdtSign(data, model.AuthToken())

		// 再次序列化
		jsonBody, _ = json.Marshal(body)
		var client = http.Client{Timeout: time.Second * 5}
		var req, err2 = http.NewRequest("POST", o.NotifyUrl, strings.NewReader(string(jsonBody)))
		if err2 != nil {
			db.Rollback()

			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Powered-By", "https://github.com/v03413/BEpusdt")
		resp, err := client.Do(req)
		if err != nil {
			db.Rollback()

			return err
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			db.Rollback()

			return fmt.Errorf("resp.StatusCode != 200")
		}

		all, _ := io.ReadAll(resp.Body)

		log.Info(fmt.Sprintf("订单回调成功[%d]：%s %s", order.Status, o.TradeId, string(all)))

		db.Commit()

		return nil
	}
	go func() {
		if err := todo(); err != nil {
			log.Warn("notify BEpusdt Error:", err.Error())
		}
	}()
}

func markNotifyFail(order model.Order, reason string) {
	log.Warn(fmt.Sprintf("订单回调失败(%v)：%s %v", order.TradeId, reason, order.SetNotifyState(model.OrderNotifyStateFail)))

	notifier.NotifyFail(order, reason)
}
