package admin

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/notifier"
	"github.com/v03413/bepusdt/app/utils"
)

type Conf struct {
}

type confReq struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

type confGetsReq struct {
	Keys []string `json:"keys" binding:"required"`
}

type confSetsReq []confReq

type notifierConf struct {
	Channel string          `json:"channel" binding:"required"`
	Params  json.RawMessage `json:"params" binding:"required"`
}

func (Conf) Set(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.SetK(model.ConfKey(strings.TrimSpace(req.Key)), strings.TrimSpace(req.Value))

	base.Ok(ctx, "配置成功")
}

func (Conf) Get(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	base.Ok(ctx, gin.H{"key": req.Key, "value": model.GetK(model.ConfKey(req.Key))})
}

func (Conf) Del(ctx *gin.Context) {
	var req confReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.Db.Where("`k` = ?", req.Key).Delete(&model.Conf{})

	base.Ok(ctx, "删除成功")
}

func (Conf) Gets(ctx *gin.Context) {
	var req confGetsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var items = make([]model.Conf, 0)
	model.Db.Where("`k` IN ?", req.Keys).Find(&items)

	var data = gin.H{}
	for _, item := range items {
		data[string(item.K)] = item.V
	}

	base.Ok(ctx, data)
}

func (Conf) Sets(ctx *gin.Context) {
	var req confSetsReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	keys := make([]string, 0)
	data := make([]model.Conf, 0)
	for _, item := range req {
		var k = strings.TrimSpace(item.Key)
		var v = strings.TrimSpace(item.Value)
		keys = append(keys, k)
		data = append(data, model.Conf{K: model.ConfKey(k), V: v})
	}

	for _, v := range data {
		if v.K == model.PaymentStaticPath && v.V != "" && !utils.IsExist(v.V) {
			base.BadRequest(ctx, "静态资源路径不存在，请确认后重新配置："+v.V)

			return
		}
	}

	model.Db.Where("`k` IN ?", keys).Delete(&model.Conf{})
	model.Db.Create(&data)

	base.Ok(ctx, "配置成功")
}

func (Conf) Notifier(ctx *gin.Context) {
	var req notifierConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var keys = []string{string(model.NotifierChannel), string(model.NotifierParams)}
	model.Db.Where("`k` IN ?", keys).Delete(&model.Conf{})
	model.Db.Create(&[]model.Conf{
		{K: model.NotifierChannel, V: req.Channel},
		{K: model.NotifierParams, V: string(req.Params)},
	})

	base.Ok(ctx, "配置成功")
}

func (Conf) NotifierTest(ctx *gin.Context) {
	err := notifier.Test()
	if err != nil {
		base.Ok(ctx, "发送测试失败："+err.Error())

		return
	}

	base.Ok(ctx, "发送测试成功")
}
