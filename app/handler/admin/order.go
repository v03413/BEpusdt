package admin

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/task/notify"
)

type Order struct {
}

type oListReq struct {
	base.ListRequest
	Name      string `json:"name"`
	Money     string `json:"money"`
	Amount    string `json:"amount"`
	OrderId   string `json:"order_id"`
	TradeId   string `json:"trade_id"`
	Status    *uint8 `json:"status"`
	Address   string `json:"address"`
	TradeType string `json:"trade_type"`
	StartAt   string `json:"start_at"`
	EndAt     string `json:"end_at"`
}

type paidReq struct {
	base.IDRequest
	RefHash string `json:"ref_hash"`
}

func (Order) List(ctx *gin.Context) {
	var req oListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	type order struct {
		model.Order
		Wallet model.Wallet `gorm:"foreignKey:MatchAddr,TradeType;references:Address,TradeType" json:"wallet"`
	}

	var data []order
	var db = model.Db

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Money != "" {
		db = db.Where("money = ?", req.Money)
	}
	if req.Amount != "" {
		db = db.Where("amount = ?", req.Amount)
	}
	if req.OrderId != "" {
		db = db.Where("order_id like ?", "%"+req.OrderId+"%")
	}
	if req.TradeId != "" {
		db = db.Where("trade_id like ?", "%"+req.TradeId+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Address != "" {
		db = db.Where("address like ?", "%"+req.Address+"%")
	}
	if req.TradeType != "" {
		db = db.Where("trade_type = ?", req.TradeType)
	}
	if req.StartAt != "" {
		db = db.Where("created_at >= ?", req.StartAt)
	}
	if req.EndAt != "" {
		db = db.Where("created_at <= ?", req.EndAt)
	}
	db = db.Where("status <> ?", 0)

	var total int64

	db.Model(&model.Order{}).Count(&total)

	err := db.Preload("Wallet").Limit(req.Size).Offset((req.Page - 1) * req.Size).Order("id " + req.Sort).Find(&data).Error
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

	type detail struct {
		model.Order
		TxUrl string `gorm:"-" json:"tx_url"`
	}

	var o model.Order
	model.Db.Where("id = ?", req.ID).Find(&o)
	if o.ID == 0 {
		base.BadRequest(ctx, "订单不存在")

		return
	}

	base.Ok(ctx, detail{
		Order: o,
		TxUrl: o.GetTxUrl(),
	})
}

func (Order) Paid(ctx *gin.Context) {
	var req paidReq
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

	var update = map[string]interface{}{
		"ref_hash":     req.RefHash,
		"status":       model.OrderStatusSuccess,
		"confirmed_at": model.Datetime(time.Now()),
	}

	err := model.Db.Model(&order).Updates(update).Error
	if err != nil {
		base.Error(ctx, err)

		return
	}

	go notify.Handle(order)

	base.Ok(ctx, "操作成功")
}

func (Order) ManualNotify(ctx *gin.Context) {
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

	if order.Status != model.OrderStatusSuccess {
		base.BadRequest(ctx, "订单状态不是交易成功,无法手动回调")

		return
	}

	go notify.Handle(order)

	base.Ok(ctx, "回调已触发")
}

func (Order) Del(ctx *gin.Context) {
	var req paidReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	err := model.Db.Where("id = ?", req.ID).Delete(&model.Order{}).Error
	if err != nil {
		base.Error(ctx, err)

		return
	}

	base.Ok(ctx, "删除成功")
}
func (Order) BatchDel(ctx *gin.Context) {
	var req base.IDListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	err := model.Db.Delete(&model.Order{}, req.IDList).Error
	if err != nil {
		base.Error(ctx, err)

		return
	}

	base.Ok(ctx, "批量删除成功")
}
