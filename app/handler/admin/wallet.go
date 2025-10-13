package admin

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
)

type Wallet struct {
}

type wAddReq struct {
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	Address   string `json:"address" binding:"required"`
	TradeType string `json:"trade_type" binding:"required"`
}

type wModReq struct {
	base.IDRequest
	Name        *string `json:"name"`
	Status      *uint8  `json:"status"`
	Address     *string `json:"address"`
	Remark      *string `json:"remark"`
	TradeType   *string `json:"trade_type"`
	OtherNotify *uint8  `json:"other_notify"`
}

type wListReq struct {
	base.ListRequest
	Name    string `json:"name"`
	Address string `json:"address"`
	Trade   string `json:"trade_type"`
}

func (Wallet) Add(ctx *gin.Context) {
	var req wAddReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	if !utils.InStrings(req.TradeType, model.SupportTradeTypes) {
		base.BadRequest(ctx, fmt.Sprintf("不支持的交易类型: %s", req.TradeType))

		return
	}

	var wallet = model.Wallet{
		Name:        strings.TrimSpace(req.Name),
		Remark:      req.Remark,
		Address:     strings.TrimSpace(req.Address),
		TradeType:   req.TradeType,
		Status:      model.WaStatusEnable,
		OtherNotify: model.WaOtherDisable,
	}

	if !wallet.IsValid() {
		base.BadRequest(ctx, "钱包地址格式不合法，请检查")

		return
	}

	if err := model.Db.Create(&wallet).Error; err != nil {
		base.Error(ctx, err)

		return
	}

	base.Response(ctx, 200, "success")
}

func (Wallet) List(ctx *gin.Context) {
	var req wListReq
	if err := ctx.ShouldBind(&req); err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	var data []model.Wallet
	var db = model.Db

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Address != "" {
		db = db.Where("address LIKE ?", "%"+req.Address+"%")
	}
	if req.Trade != "" {
		db = db.Where("trade_type LIKE ?", "%"+req.Trade+"%")
	}

	var total int64

	db.Model(&model.Wallet{}).Count(&total)

	err := db.Limit(req.Size).Offset((req.Page - 1) * req.Size).Order("id " + req.Sort).Find(&data).Error
	if err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	base.Response(ctx, 200, data, total)
}

func (Wallet) Mod(ctx *gin.Context) {
	var req wModReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var w model.Wallet
	model.Db.Where("id = ?", req.ID).Find(&w)
	if w.ID == 0 {
		base.BadRequest(ctx, "钱包不存在")

		return
	}

	if req.Name != nil {
		w.Name = strings.TrimSpace(*req.Name)
	}
	if req.Remark != nil {
		w.Remark = *req.Remark
	}
	if req.Address != nil {
		w.Address = strings.TrimSpace(*req.Address)
	}
	if req.TradeType != nil {
		if !utils.InStrings(*req.TradeType, model.SupportTradeTypes) {
			base.BadRequest(ctx, fmt.Sprintf("不支持的交易类型: %s", *req.TradeType))

			return
		}

		w.TradeType = *req.TradeType
	}
	if req.Status != nil {
		w.Status = *req.Status
	}
	if req.OtherNotify != nil {
		w.OtherNotify = *req.OtherNotify
	}

	if !w.IsValid() {
		base.BadRequest(ctx, "钱包地址格式不合法，请检查")

		return
	}

	if err := model.Db.Save(&w).Error; err != nil {
		base.Error(ctx, err)

		return
	}

	base.Response(ctx, 200, "修改成功")
}

func (Wallet) Del(ctx *gin.Context) {
	var req base.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	model.Db.Where("id = ?", req.ID).Delete(&model.Wallet{})

	base.Response(ctx, 200, "删除成功")
}
