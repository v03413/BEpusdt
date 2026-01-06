package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/auth"
)

func authInit(e *gin.Engine) {
	var authRtr = e.Group("/api/auth")
	var authHdr = new(auth.Auth)
	{
		GetRegister(authRtr, "/info", true, authHdr.Info)
		GetRegister(authRtr, "/menu", true, authHdr.Menu)
		PostRegister(authRtr, "/login", false, authHdr.Login)
		PostRegister(authRtr, "/logout", true, authHdr.Logout)
		PostRegister(authRtr, "/set_password", true, authHdr.SetPassword)
	}
}
