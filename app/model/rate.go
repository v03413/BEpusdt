package model

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/utils"
	"gorm.io/gorm"
)

const (
	FiatCNY = "CNY"
	FiatUSD = "USD"
	FiatJPY = "JPY"
	FiatEUR = "EUR"
	FiatGBP = "GBP"
)

// SupportTradeFiat 支持的法币交易类型
var SupportTradeFiat = []string{FiatCNY, FiatUSD, FiatJPY, FiatEUR, FiatGBP}

var SupportTradeCrypto = []string{string(TokenTypeUSDT), string(TokenTypeUSDC), string(TokenTypeTRX)}

type Rate struct {
	Id
	Rate    string  `Gorm:"column:rate;type:varchar(32);not null;comment:订单汇率" json:"rate"`
	Fiat    string  `Gorm:"column:fiat;type:varchar(16);not null;comment:法币" json:"fiat"`
	Crypto  string  `Gorm:"column:crypto;type:varchar(16);not null;comment:加密货币" json:"crypto"`
	RawRate float64 `Gorm:"column:raw_rate;type:decimal(10,4);not null;comment:基准汇率" json:"raw_rate"`
	Syntax  string  `Gorm:"column:syntax;type:varchar(32);not null;default:'';comment:浮动语法" json:"syntax"`
	AutoTimeAt
}

func (r *Rate) TableName() string {
	return "bep_rate"
}

func (r *Rate) BeforeCreate(tx *gorm.DB) error {
	var syntax = GetK(ConfKey(fmt.Sprintf("rate_float_%s_%s", r.Crypto, r.Fiat)))
	if syntax == "" {

		return nil
	}

	r.Syntax = syntax
	r.Rate = cast.ToString(ParseFloatRate(syntax, cast.ToFloat64(r.RawRate)))

	return nil
}

func CoingeckoRate() {
	var c = &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=tron,tether,usd-coin&vs_currencies=%s", strings.Join(SupportTradeFiat, ",")))
	if err != nil {
		log.Warn("汇率同步错误", err)

		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("汇率同步错误", err)

		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Warn("汇率同步错误 StatusCode != 200", resp.StatusCode)

		return
	}

	var rows = make([]Rate, 0)
	var tokenTypes = map[string]TokenType{
		"tron":     TokenTypeTRX,
		"tether":   TokenTypeUSDT,
		"usd-coin": TokenTypeUSDC,
	}
	for k, v := range gjson.ParseBytes(body).Map() {
		var token, ok = tokenTypes[k]
		if !ok {

			continue
		}

		for fiat, val := range v.Map() {
			rows = append(rows, Rate{
				Rate:    val.String(),
				Fiat:    strings.ToUpper(fiat),
				Crypto:  string(token),
				RawRate: val.Float(),
			})
		}
	}

	Db.Create(&rows)
	log.Info("汇率同步成功")
}

func ParseFloatRate(syntax string, rawVal float64) float64 {
	if syntax == "" {

		return rawVal
	}

	if utils.IsNumber(syntax) {

		return cast.ToFloat64(syntax)
	}

	match, err := regexp.MatchString(`^[~+-]\d+(\.\d+)?$`, syntax)
	if !match || err != nil {
		log.Error("浮动语法解析错误", err)

		return 0
	}

	var act = syntax[0:1]
	var raw = decimal.NewFromFloat(rawVal)
	var base = decimal.NewFromFloat(cast.ToFloat64(syntax[1:]))
	var result float64 = 0

	switch act {
	case "~":
		result = raw.Mul(base).InexactFloat64()
	case "+":
		result = raw.Add(base).InexactFloat64()
	case "-":
		result = raw.Sub(base).InexactFloat64()
	}

	return round(result, 4) // 保留4位小数
}

func round(val float64, precision int) float64 {
	// Round 四舍五入，ROUND_HALF_UP 模式实现
	// 返回将 val 根据指定精度 precision（十进制小数点后数字的数目）进行四舍五入的结果。precision 也可以是负数或零。

	if precision == 0 {
		return math.Round(val)
	}

	p := math.Pow10(precision)
	if precision < 0 {
		return math.Floor(val*p+0.5) * math.Pow10(-precision)
	}

	return math.Floor(val*p+0.5) / p
}

func getOrderRate(token TokenType, fiat, syntax string) (decimal.Decimal, error) {
	var r Rate
	Db.Where("crypto = ? and fiat = ?", token, fiat).Order("created_at desc").Limit(1).Find(&r)
	if r.ID == 0 {

		return decimal.Decimal{}, fmt.Errorf("创建失败，请检查汇率同步是否正常：%s %s", token, fiat)
	}

	if syntax == "" {

		return decimal.NewFromString(r.Rate)
	}

	return decimal.NewFromFloat(ParseFloatRate(syntax, r.RawRate)), nil
}
