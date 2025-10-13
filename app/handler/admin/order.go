package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
)

type Order struct {
}

type oListReq struct {
	base.ListRequest
	Name      *string `json:"name"`
	Money     *string `json:"money"`
	Amount    *string `json:"amount"`
	OrderId   *string `json:"order_id"`
	TradeId   *string `json:"trade_id"`
	Status    *uint8  `json:"status"`
	Address   *string `json:"address"`
	TradeType *string `json:"trade_type"`
}

func (Order) List(ctx *gin.Context) {
	var req oListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	var data []model.Order
	var db = model.Db

	if req.Name != nil {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Money != nil {
		db = db.Where("money = ?", *req.Money)
	}
	if req.Amount != nil {
		db = db.Where("amount = ?", *req.Amount)
	}
	if req.OrderId != nil {
		db = db.Where("order_id like ?", "%"+*req.OrderId+"%")
	}
	if req.TradeId != nil {
		db = db.Where("trade_id like ?", "%"+*req.TradeId+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Address != nil {
		db = db.Where("address like ?", "%"+*req.Address+"%")
	}
	if req.TradeType != nil {
		db = db.Where("trade_type = ?", *req.TradeType)
	}

	var total int64

	db.Model(&model.Order{}).Count(&total)

	err := db.Limit(req.Size).Offset((req.Page - 1) * req.Size).Order("id " + req.Sort).Find(&data).Error
	if err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	base.Response(ctx, 200, data, total)
}

func (Order) Detail(ctx *gin.Context) {
	var req base.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var order model.Order
	model.Db.Where("id = ?", req.ID).Find(&order)
	if order.ID == 0 {
		base.BadRequest(ctx, "订单不存在")

		return
	}

	base.Ok(ctx, order)
}

func (Order) Paid(ctx *gin.Context) {
	var req base.IDRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var order model.Order
	model.Db.Where("id = ?", req.ID).Find(&order)
	if order.ID == 0 {
		base.BadRequest(ctx, "订单不存在")

		return
	}

	err := model.Db.Model(&order).Update("status", model.OrderStatusSuccess).Error
	if err != nil {
		base.Error(ctx, err)

		return
	}

	base.Ok(ctx, "操作成功")
}
