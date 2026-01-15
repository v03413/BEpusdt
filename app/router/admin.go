package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/admin"
)

func adminInit(e *gin.Engine) {
	var confRtr = e.Group("/api/conf")
	var confHdr = new(admin.Conf)
	{
		PostRegister(confRtr, "/set", true, confHdr.Set)
		PostRegister(confRtr, "/get", true, confHdr.Get)
		PostRegister(confRtr, "/del", true, confHdr.Del)
		PostRegister(confRtr, "/gets", true, confHdr.Gets)
		PostRegister(confRtr, "/sets", true, confHdr.Sets)
		PostRegister(confRtr, "/notifier", true, confHdr.Notifier)
		PostRegister(confRtr, "/notifier_test", true, confHdr.NotifierTest)
		PostRegister(confRtr, "/reset_api_auth_token", true, confHdr.ResetApiAuthToken)
	}

	var walletRtr = e.Group("/api/wallet")
	var walletHdr = new(admin.Wallet)
	{
		PostRegister(walletRtr, "/add", true, walletHdr.Add)
		PostRegister(walletRtr, "/list", true, walletHdr.List)
		PostRegister(walletRtr, "/mod", true, walletHdr.Mod)
		PostRegister(walletRtr, "/del", true, walletHdr.Del)
	}

	var orderRtr = e.Group("/api/order")
	var orderHdr = new(admin.Order)
	{
		PostRegister(orderRtr, "/list", true, orderHdr.List)
		PostRegister(orderRtr, "/detail", true, orderHdr.Detail)
		PostRegister(orderRtr, "/paid", true, orderHdr.Paid)
		PostRegister(orderRtr, "/manual_notify", true, orderHdr.ManualNotify)
	}

	var rateRtr = e.Group("/api/rate")
	var rateHdr = new(admin.Rate)
	{
		PostRegister(rateRtr, "/list", true, rateHdr.List)
		PostRegister(rateRtr, "/syntax", true, rateHdr.Syntax)
		PostRegister(rateRtr, "/set_syntax", true, rateHdr.SetSyntax)
		PostRegister(rateRtr, "/sync", true, rateHdr.Sync)
	}

	var dashboardRtr = e.Group("/api/dashboard")
	var dashboardHdr = new(admin.Dashboard)
	{
		PostRegister(dashboardRtr, "/home", true, dashboardHdr.Home)
	}
}
