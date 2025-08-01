package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/web/epay"
	"net/http"
)

// epaySubmit 【兼容】易支付提交
func epaySubmit(ctx *gin.Context) {
	var rawParams = ctx.Request.URL.Query()
	if ctx.Request.Method == http.MethodPost {
		if err := ctx.Request.ParseForm(); err != nil {
			ctx.String(200, "参数解析错误："+err.Error())

			return
		}

		rawParams = ctx.Request.PostForm
	}

	var data = make(map[string]string)
	for k, v := range rawParams {
		if len(v) == 0 {
			data[k] = ""

			continue
		}

		data[k] = v[0]
	}

	if data["pid"] != epay.Pid {
		ctx.String(200, "BEpusdt 易支付兼容模式，商户号【PID】必须固定为"+epay.Pid)

		return
	}

	if epay.Sign(data, conf.GetAuthToken()) != data["sign"] {
		ctx.String(200, "签名错误")

		return
	}

	var tradeType = model.OrderTradeTypeUsdtTrc20
	if v, ok := data["type"]; ok {

		tradeType = cast.ToString(v)
	}

	var params = orderParams{
		Money:       cast.ToFloat64(data["money"]),
		ApiType:     model.OrderApiTypeEpay,
		PayAddress:  "",
		OrderId:     data["out_trade_no"],
		TradeType:   tradeType,
		RedirectUrl: data["return_url"],
		NotifyUrl:   data["notify_url"],
		Name:        data["name"],
	}

	var order, err = buildOrder(params)
	if err != nil {
		ctx.String(200, fmt.Sprintf("订单创建失败：%v", err))

		return
	}

	// 解析请求地址
	var host = "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("%s/pay/checkout-counter/%s", conf.GetAppUri(host), order.TradeId))
}
