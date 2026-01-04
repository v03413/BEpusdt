package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
)

type Rate struct {
}

type rateListReq struct {
	base.ListRequest
	Fiat     string   `json:"fiat"`
	Crypto   string   `json:"crypto"`
	Datetime []string `json:"datetime"`
}

type rateSetSyntaxReq struct {
	Fiat   string `json:"fiat"`
	Crypto string `json:"crypto"`
	Syntax string `json:"syntax"`
}

func (Rate) List(ctx *gin.Context) {
	var req rateListReq
	if err := ctx.ShouldBind(&req); err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	var data []model.Rate
	var db = model.Db
	if req.Fiat != "" {
		db = db.Where("fiat = ?", req.Fiat)
	}
	if req.Crypto != "" {
		db = db.Where("crypto = ?", req.Crypto)
	}
	if len(req.Datetime) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", req.Datetime[0], req.Datetime[1])
	}

	var total int64

	db.Model(&model.Rate{}).Count(&total)

	err := db.Limit(req.Size).Offset((req.Page - 1) * req.Size).Order("id " + req.Sort).Find(&data).Error
	if err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	base.Response(ctx, 200, data, total)
}

func (Rate) Syntax(ctx *gin.Context) {
	type Syntax struct {
		Key    string `json:"key"`
		Fiat   string `json:"fiat"`
		Crypto string `json:"crypto"`
		Syntax string `json:"syntax"`
	}

	var data = make([]Syntax, 0)

	for fiat := range model.SupportFiat {
		for token := range model.SupportCrypto {
			var k = fmt.Sprintf("rate_float_%s_%s", token, fiat)
			data = append(data, Syntax{
				Key:    k,
				Fiat:   string(fiat),
				Crypto: string(token),
				Syntax: model.GetK(model.ConfKey(k)),
			})
		}
	}

	base.Ok(ctx, data)
}

func (Rate) SetSyntax(ctx *gin.Context) {
	var req rateSetSyntaxReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.SetK(model.ConfKey("rate_float_"+req.Crypto+"_"+req.Fiat), req.Syntax)

	base.Ok(ctx, "设置成功")
}

func (Rate) Sync(ctx *gin.Context) {
	model.CoingeckoRate()

	base.Ok(ctx, "同步完成。")
}
