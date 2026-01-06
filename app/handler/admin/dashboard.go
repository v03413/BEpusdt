package admin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
)

type Dashboard struct {
}

type homeReq struct {
	Fiat string `json:"fiat" binding:"required"`
}

func (Dashboard) Home(ctx *gin.Context) {
	var req homeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var rows = make([]model.Order, 0)
	model.Db.Where("fiat = ? and status = ?", req.Fiat, model.OrderStatusSuccess).Find(&rows)

	var today = time.Now().Format("2006-01-02")
	var totalMoney, todayMoney float64
	var totalCount, todayCount int64
	var monthly = make(map[string]float64)
	var tokenMap = map[model.Crypto]float64{
		model.USDT: 0,
		model.USDC: 0,
		model.TRX:  0,
		model.BNB:  0,
		model.ETH:  0,
	}

	for _, itm := range rows {
		money := cast.ToFloat64(itm.Money)

		totalMoney += money
		totalCount++
		if itm.CreatedAt.Format("2006-01-02") == today {
			todayMoney += money
			todayCount++
		}

		if crypto, err := model.GetCrypto(itm.TradeType); err == nil {
			if _, ok := tokenMap[crypto]; !ok {
				tokenMap[crypto] = 0
			}

			tokenMap[crypto] += money
		}

		var month = itm.CreatedAt.Time().Format("2006/01")
		if _, ok := monthly[month]; !ok {
			monthly[month] = 0
		}
		monthly[month] += money
	}

	base.Ok(ctx, gin.H{
		"total_money": totalMoney,
		"total_count": totalCount,
		"today_money": todayMoney,
		"today_count": todayCount,
		"token_map":   tokenMap,
		"monthly":     monthly,
	})
}
