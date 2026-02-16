package epusdt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/handler/epay"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

type Epusdt struct{}

type createReq struct {
	Amount      float64    `json:"amount" binding:"required"`
	OrderID     string     `json:"order_id" binding:"required"`
	NotifyURL   string     `json:"notify_url" binding:"required"`
	RedirectURL string     `json:"redirect_url" binding:"required"`
	Signature   string     `json:"signature" binding:"required"`
	Name        string     `json:"name"`
	Fiat        model.Fiat `json:"fiat"`
	TradeType   string     `json:"trade_type"`
	Address     string     `json:"address"`
	Timeout     int64      `json:"timeout"`
	Rate        string     `json:"rate"`
}

type cancelReq struct {
	TradeID   string `json:"trade_id" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type methodsReq struct {
	TradeID  string `json:"trade_id" binding:"required"`
	Currency string `json:"currency"`
}

type methodItem struct {
	Amount          string `json:"amount"`
	ActualAmount    string `json:"actual_amount"`
	Fiat            string `json:"fiat"`
	ExchangeRate    string `json:"exchange_rate"`
	Currency        string `json:"currency"`
	Network         string `json:"network"`
	TokenNetName    string `json:"token_net_name"`
	TokenCustomName string `json:"token_custom_name"`
	IsPopular       bool   `json:"is_popular"`
}

func (Epusdt) SignVerify(ctx *gin.Context) {
	rawData, err := ctx.GetRawData()
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("json 数据读取错误 %s", err.Error())))
		ctx.Abort()

		return
	}

	m := make(map[string]any)
	if err = json.Unmarshal(rawData, &m); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("json 数据解析错误 %s", err.Error())))
		ctx.Abort()

		return
	}

	sign, ok := m["signature"]
	if !ok {
		ctx.JSON(200, respFailJson("签名丢失"))
		ctx.Abort()

		return
	}

	if utils.EpusdtSign(m, model.AuthToken()) != sign {
		ctx.JSON(200, respFailJson("签名错误"))
		ctx.Abort()

		return
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawData)) // 回写数据
	ctx.Next()
}

func (Epusdt) CreateTransaction(ctx *gin.Context) {
	var req createReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("请求参数错误：%s", err.Error())))

		return
	}

	// 解析请求地址
	host := "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	if req.Fiat == "" {
		req.Fiat = model.CNY
	}
	if req.TradeType == "" {
		req.TradeType = string(model.UsdtTrc20)
	}

	order, err := model.BuildOrder(model.OrderParams{
		Money:       decimal.NewFromFloat(req.Amount),
		ApiType:     model.OrderApiTypeEpusdt,
		Address:     req.Address,
		OrderId:     req.OrderID,
		TradeType:   model.TradeType(req.TradeType),
		RedirectUrl: req.RedirectURL,
		NotifyUrl:   req.NotifyURL,
		Name:        req.Name,
		Timeout:     req.Timeout,
		Rate:        req.Rate,
		Fiat:        req.Fiat,
	})
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("订单创建失败：%s", err.Error())))

		return
	}

	log.Info(fmt.Sprintf("订单创建成功 商户订单：%s", req.OrderID))

	// 返回响应数据
	ctx.JSON(200, respSuccJson(gin.H{
		"fiat":            order.Fiat,
		"trade_type":      order.TradeType,
		"trade_id":        order.TradeId,
		"order_id":        order.OrderId,
		"status":          order.Status,
		"amount":          order.Money,
		"actual_amount":   order.Amount,
		"token":           order.Address,
		"expiration_time": uint64(order.ExpiredAt.Sub(time.Now()).Seconds()),
		"payment_url":     model.CheckoutCounter(host, order.TradeId),
	}))
}

func (Epusdt) CancelTransaction(ctx *gin.Context) {
	var req cancelReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("请求参数错误：%s", err.Error())))

		return
	}

	order, ok := model.GetTradeOrder(req.TradeID)
	if !ok {
		ctx.JSON(200, respFailJson("订单不存在"))

		return
	}

	if order.Status != model.OrderStatusWaiting {
		ctx.JSON(200, respFailJson(fmt.Sprintf("当前订单(%s)状态不允许取消", req.TradeID)))

		return
	}

	if err := order.SetCanceled(); err != nil {
		ctx.JSON(200, respFailJson(fmt.Sprintf("订单取消失败：%s", err.Error())))

		return
	}

	ctx.JSON(200, respSuccJson(gin.H{"trade_id": req.TradeID}))
}

func (Epusdt) CheckoutCounter(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.String(200, "订单不存在")

		return
	}

	uri, err := url.ParseRequestURI(order.ReturnUrl)
	if err != nil {
		ctx.String(200, "同步地址错误")
		log.Error("同步地址解析错误", err.Error())

		return
	}

	ctx.HTML(200, string(order.TradeType+".html"), gin.H{
		"http_host":  uri.Host,
		"amount":     order.Amount,
		"address":    order.Address,
		"expire":     int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
		"return_url": order.ReturnUrl,
		"usdt_rate":  order.Rate,
		"trade_id":   tradeId,
		"order_id":   order.OrderId,
		"trade_type": order.TradeType,
		"money":      order.Money,
		"fiat":       order.Fiat,
	})
}

func (Epusdt) CheckStatus(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.JSON(200, respFailJson("订单不存在"))

		return
	}

	var returnUrl string
	if order.Status == model.OrderStatusSuccess {
		returnUrl = order.ReturnUrl
		if order.ApiType == model.OrderApiTypeEpay {
			// 易支付兼容
			returnUrl = fmt.Sprintf("%s?%s", returnUrl, epay.BuildNotifyParams(order))
		}
	}

	// 返回响应数据
	ctx.JSON(200, gin.H{
		"trade_id":   tradeId,
		"trade_hash": order.RefHash,
		"status":     order.Status,
		"return_url": returnUrl,
	})
}

// Methods API
func (Epusdt) GetPaymentMethods(ctx *gin.Context) {
	var req methodsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(500, respFailJson(fmt.Sprintf("request error: %s", err.Error())))
		return
	}

	order, ok := model.GetTradeOrder(req.TradeID)
	if !ok {
		ctx.JSON(400, respFailJson("order not found"))
		return
	}

	if order.Status == model.OrderStatusExpired {
		ctx.JSON(400, respFailJson("order expired"))
		return
	}

	// 强制使用订单的法币类型
	fiat := order.Fiat

	var methods []methodItem
	allTrades := model.GetAllTradeConfig()

	for tradeTypeStr, conf := range allTrades {
		// 如果指定了货币，则进行过滤
		if req.Currency != "" && string(conf.Crypto) != req.Currency {
			continue
		}

		// 检查是否有可用钱包
		var count int
		tradeType := model.TradeType(tradeTypeStr)
		count = len(model.GetAvailableAddress(tradeType))

		if count == 0 {
			continue
		}

		// 获取汇率配置的浮动语法
		syntax := model.GetK(model.ConfKey(fmt.Sprintf("rate_float_%s_%s", conf.Crypto, fiat)))

		// 获取汇率
		rate, err := model.GetOrderRate(conf.Crypto, fiat, syntax)
		if err != nil {
			log.Error(fmt.Sprintf("GetPaymentMethods: get order rate error: %s", err.Error()))
			continue
		}

		// 计算实际支付金额 (加密货币)
		// Money 是法币金额
		moneyDecimal, _ := decimal.NewFromString(order.Money)

		// 计算精度
		atom, precision := model.GetAtomicity(model.TradeType(tradeTypeStr))
		actualAmount := moneyDecimal.DivRound(rate, precision)
		if actualAmount.LessThan(atom) {
			actualAmount = atom
		}
		actualAmountStr := actualAmount.String()

		methods = append(methods, methodItem{
			Amount:          order.Money,
			ActualAmount:    actualAmountStr,
			Fiat:            string(fiat),
			ExchangeRate:    rate.String(),
			Currency:        string(conf.Crypto),
			Network:         string(conf.Network),
			TokenNetName:    string(conf.NetworkName),
			TokenCustomName: "",    // 暂为空
			IsPopular:       false, // 暂为 false
		})
	}

	// Sort by Currency A-Z
	sort.Slice(methods, func(i, j int) bool {
		if methods[i].Currency != methods[j].Currency {
			return methods[i].Currency < methods[j].Currency
		}
		return methods[i].Network < methods[j].Network
	})

	ctx.JSON(200, respSuccJson(gin.H{
		"methods": methods,
	}))
}

func respFailJson(message string) gin.H {

	return gin.H{"status_code": 400, "message": message}
}

func respSuccJson(data interface{}) gin.H {

	return gin.H{"status_code": 200, "message": "success", "data": data, "request_id": ""}
}
