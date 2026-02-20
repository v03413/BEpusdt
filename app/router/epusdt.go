package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/epusdt"
)

func epusdtInit(engine *gin.Engine) {
	epGrp := engine.Group("/pay")
	epHdr := new(epusdt.Epusdt)
	{
		epGrp.GET("/checkout-counter/:trade_id", epHdr.CheckoutCounter)
		epGrp.GET("/order/:trade_id", epHdr.CheckoutCashier)
		epGrp.GET("/check-status/:trade_id", epHdr.CheckStatus)
	}

	orderGrp := engine.Group("/api/v1/order")
	{
		orderGrp.Use(epHdr.SignVerify)
		orderGrp.POST("/create-transaction", epHdr.CreateTransaction)
		orderGrp.POST("/cancel-transaction", epHdr.CancelTransaction)
		orderGrp.POST("/create-order", epHdr.CreateOrder)
	}

	payApiGrp := engine.Group("/api/v1/pay")
	{
		payApiGrp.POST("/methods", epHdr.GetPaymentMethods)
		payApiGrp.POST("/update-order", epHdr.UpdateOrder)
		payApiGrp.POST("/notify", epHdr.Notify)
	}
}
