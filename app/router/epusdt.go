package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/epusdt"
)

func epusdtInit(engine *gin.Engine) {
	epGrp := engine.Group("/pay")
	epHdr := new(epusdt.Epusdt)
	{
		epGrp.GET("/checkout/:trade_id", epHdr.Checkout)
	}

	orderGrp := engine.Group("/api/v1/order")
	{
		orderGrp.Use(epHdr.SignVerify)
		orderGrp.POST("/create-transaction", epHdr.CreateTransaction)
		orderGrp.POST("/cancel-transaction", epHdr.CancelTransaction)
		orderGrp.POST("/create-order", epHdr.CreateOrder)
	}

	payGrp := engine.Group("/api/v1/pay")
	{
		payGrp.POST("/info", epHdr.Info)
		payGrp.POST("/notify", epHdr.Notify)
		payGrp.POST("/methods", epHdr.GetMethods)
		payGrp.POST("/update-order", epHdr.UpdateOrder)
	}
}
