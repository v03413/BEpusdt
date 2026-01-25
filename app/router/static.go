package router

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/v03413/bepusdt/static"
)

func staticInit(e *gin.Engine) {
	customPath := model.GetK(model.PaymentStaticPath)
	if customPath != "" && utils.IsExist(customPath) {
		initCustomPayment(e, customPath)

		return
	}

	initDefaultPayment(e)
}

func initCustomPayment(e *gin.Engine, path string) {
	tmpl := template.New("customer")

	template.Must(tmpl.ParseGlob(filepath.Join(path, "views", "*.html")))
	template.Must(tmpl.ParseFS(static.Secure, "secure/secure.html"))
	e.SetHTMLTemplate(tmpl)

	e.StaticFS("/payment/assets", http.Dir(filepath.Join(path, "assets")))
	e.StaticFS("/secure/assets", http.FS(subFS(static.Secure, "secure/assets")))

	log.Info("成功注册自定义静态资源路径：", path)
}

func initDefaultPayment(e *gin.Engine) {
	tmpl := template.New("default")

	template.Must(tmpl.ParseFS(static.Secure, "secure/secure.html"))
	template.Must(tmpl.ParseFS(static.Payment, "payment/views/*.html"))
	e.SetHTMLTemplate(tmpl)

	e.StaticFS("/payment/assets", http.FS(subFS(static.Payment, "payment/assets")))
	e.StaticFS("/secure/assets", http.FS(subFS(static.Secure, "secure/assets")))
}

func subFS(src fs.FS, dir string) fs.FS {
	sub, _ := fs.Sub(src, dir)

	return sub
}
