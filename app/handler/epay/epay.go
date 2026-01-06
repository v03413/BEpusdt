package epay

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
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
	Money      float64    `form:"money" json:"money" binding:"required,gt=0"`
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
	if ctx.Request.Method != http.MethodPost {
		err = ctx.ShouldBindQuery(&data)
	} else {
		err = ctx.ShouldBind(&data)
	}

	if err != nil {
		ctx.String(200, "参数解析错误："+err.Error())

		return
	}
	if data.Pid != Pid {
		ctx.String(200, "BEpusdt 易支付兼容模式，商户号【PID】必须固定为"+Pid)

		return
	}

	if e.sign(data, model.AuthToken()) != data.Sign {
		ctx.String(200, "签名错误")

		return
	}

	e.fillDefaultParams(&data)

	money := decimal.NewFromFloat(data.Money)
	if money.LessThanOrEqual(decimal.Zero) {
		ctx.String(200, "参数 money 格式错误，必须是一个大于0的数字")

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
		ctx.String(200, fmt.Sprintf("订单创建失败：%v", err))

		return
	}

	// 解析请求地址
	var host = "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	ctx.Redirect(http.StatusFound, model.CheckoutCounter(host, order.TradeId))
}

// 填充默认参数
func (e Epay) fillDefaultParams(data *submit) {
	if data.Pid == "" {
		data.Pid = Pid
	}
	if data.Type == "" {
		data.Type = string(model.UsdtTrc20)
	}
	if data.Fiat == "" {
		data.Fiat = model.CNY
	}
}

func (e Epay) sign(params submit, token string) string {
	jsonBytes, err := json.Marshal(params)
	if err != nil {

		return ""
	}

	result := gjson.ParseBytes(jsonBytes)
	fields := make(map[string]string)
	result.ForEach(func(k, value gjson.Result) bool {
		fields[k.String()] = value.String()

		return true
	})

	// 提取 keys 并排序
	var keys = make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// 构建签名字符串
	signStr := ""
	for _, k := range keys {
		if k != "sign" && k != "sign_type" && fields[k] != "" {
			signStr += fmt.Sprintf("%s=%s&", k, fields[k])
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
