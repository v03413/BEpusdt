package router

import (
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/handler/epay"
)

func epayInit(engine *gin.Engine) {
	epHdr := new(epay.Epay)
	{
		engine.POST("/submit.php", epHdr.Submit)
		engine.GET("/submit.php", epHdr.Submit)
	}
}
