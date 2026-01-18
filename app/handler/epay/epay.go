package epay

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

const Pid = "1000" // 固定商户号

type Epay struct {
}

type submit struct {
	Pid        string     `form:"pid" json:"pid" binding:"required"` // 商户号
	Type       string     `form:"type" json:"type" binding:"required"`
	NotifyURL  string     `form:"notify_url" json:"notify_url" binding:"required"`
	ReturnURL  string     `form:"return_url" json:"return_url" binding:"required"`
	OutTradeNo string     `form:"out_trade_no" json:"out_trade_no" binding:"required"`
	Name       string     `form:"name" json:"name" binding:"required"`
	Money      string     `form:"money" json:"money" binding:"required,gt=0"`
	Sign       string     `form:"sign" json:"sign" binding:"required"`
	Fiat       model.Fiat `form:"fiat" json:"fiat"`
	Rate       string     `form:"rate" json:"rate"`
	Timeout    string     `form:"timeout" json:"timeout"`
	Address    string     `form:"address" json:"address"`
}

// Submit 【兼容】易支付提交
func (e Epay) Submit(ctx *gin.Context) {
	var err error
	var data submit

	var dataMap = make(map[string]string)
	if ctx.Request.Method != http.MethodPost {
		for key, values := range ctx.Request.URL.Query() {
			if len(values) > 0 {
				dataMap[key] = values[0]
			}
		}
	} else {
		_ = ctx.Request.ParseForm()
		for key, values := range ctx.Request.PostForm {
			if len(values) > 0 {
				dataMap[key] = values[0]
			}
		}
	}

	data, err = e.verify(dataMap)
	if err != nil {
		ctx.String(200, err.Error())

		return
	}

	if e.sign(dataMap, model.AuthToken()) != data.Sign {
		ctx.String(200, "签名错误")

		return
	}

	money, err := decimal.NewFromString(data.Money)
	if err != nil {
		ctx.String(200, "参数 money 解析错误，"+err.Error())

		return
	}

	var order, err2 = model.BuildOrder(model.OrderParams{
		Money:       money,
		ApiType:     model.OrderApiTypeEpay,
		Address:     data.Address,
		OrderId:     data.OutTradeNo,
		TradeType:   model.TradeType(data.Type),
		RedirectUrl: data.ReturnURL,
		NotifyUrl:   data.NotifyURL,
		Name:        data.Name,
		Timeout:     cast.ToInt64(data.Timeout),
		Rate:        data.Rate,
		Fiat:        data.Fiat,
	})
	if err2 != nil {
		ctx.String(200, fmt.Sprintf("订单创建失败：%v", err2))

		return
	}

	// 解析请求地址
	var host = "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	ctx.Redirect(http.StatusFound, model.CheckoutCounter(host, order.TradeId))
}

// verify 验证请求参数
func (e Epay) verify(data map[string]string) (submit, error) {
	var params = submit{}

	pid, ok := data["pid"]
	if !ok || pid != Pid {

		return params, fmt.Errorf("BEpusdt 易支付兼容模式，商户号【PID】必须固定为" + Pid)
	}

	var requiredFields = []string{"pid", "type", "out_trade_no", "notify_url", "return_url", "name", "money", "sign"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok || data[field] == "" {

			return params, fmt.Errorf("参数 %s 缺失或为空", field)
		}
	}

	params.Pid = pid
	params.Type = data["type"]
	params.OutTradeNo = data["out_trade_no"]
	params.NotifyURL = data["notify_url"]
	params.ReturnURL = data["return_url"]
	params.Name = data["name"]
	params.Money = data["money"]
	params.Sign = data["sign"]

	if address, ok := data["address"]; ok {
		params.Address = address
	}
	if rate, ok := data["rate"]; ok && rate != "" {
		params.Rate = rate
	}
	if timeout, ok := data["timeout"]; ok && timeout != "" {
		params.Timeout = timeout
	}

	fiat, ok := data["fiat"]
	if ok && fiat != "" {
		params.Fiat = model.Fiat(fiat)
	} else {
		params.Fiat = model.CNY
	}

	return params, nil
}

func (e Epay) sign(data map[string]string, token string) string {
	var keys = make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	signStr := ""
	for _, k := range keys {
		if k != "sign" && k != "sign_type" && data[k] != "" {
			signStr += fmt.Sprintf("%s=%s&", k, data[k])
		}
	}
	signStr = signStr[:len(signStr)-1]
	signStr += token

	// 计算 MD5
	hash := md5.New()
	hash.Write([]byte(signStr))
	md5sum := hex.EncodeToString(hash.Sum(nil))

	return md5sum
}

func BuildNotifyParams(order model.Order) string {
	var signStr = utils.Md5String(fmt.Sprintf("money=%s&name=%s&out_trade_no=%s&pid=%s&trade_no=%s&trade_status=TRADE_SUCCESS&type=%s",
		cast.ToString(order.Money), order.Name, order.OrderId, Pid, order.TradeId, order.TradeType) + model.AuthToken())
	var params = fmt.Sprintf("money=%s&name=%s&out_trade_no=%s&pid=%s&trade_no=%s&trade_status=TRADE_SUCCESS&type=%s",
		cast.ToString(order.Money), url.QueryEscape(order.Name), url.QueryEscape(order.OrderId), Pid, order.TradeId, order.TradeType)

	return fmt.Sprintf("%s&sign=%s", params, signStr)
}
