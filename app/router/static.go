package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/static"
)

func staticInit(e *gin.Engine) {
	tmpl := template.New("default")

	template.Must(tmpl.ParseFS(static.Secure, "secure/secure.html"))
	template.Must(tmpl.ParseFS(static.Payment, "payment/views/*.html"))

	// 注册收银台模板和静态资源
	entries, err := fs.ReadDir(static.Checkout, "checkout")
	if err != nil {
		panic(errors.New(fmt.Sprintf("收银台模板读取异常：%v", err)))
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		dir := path.Join("checkout", name)
		checkout, err := readCheckoutInfoFromFS(static.Checkout, dir)
		if err != nil {
			panic(err)
		}

		if registerTemplatesFromFS(tmpl, static.Checkout, dir, name) {
			model.RegisterCheckout(name, checkout)
			log.Info("前台收银模板注册成功：", name)
		}

		registerAssetsFromFS(e, static.Checkout, dir, "/checkout/"+name+"/assets")
	}

	e.SetHTMLTemplate(tmpl)
	e.StaticFS("/payment/assets", http.FS(subFS(static.Payment, "payment/assets")))
	e.StaticFS("/secure/assets", http.FS(subFS(static.Secure, "secure/assets")))
}

func readCheckoutInfoFromFS(src fs.FS, baseDir string) (model.Checkout, error) {
	infoPath := path.Join(baseDir, "checkout.json")
	data, err := fs.ReadFile(src, infoPath)
	if err != nil {
		return model.Checkout{}, fmt.Errorf("checkout.json 不存在或读取失败: %w", err)
	}

	var checkout model.Checkout
	if err := json.Unmarshal(data, &checkout); err != nil {
		return model.Checkout{}, fmt.Errorf("解析 checkout.json 失败: %w", err)
	}

	return checkout, nil
}

func registerTemplatesFromFS(tmpl *template.Template, src fs.FS, baseDir, namePrefix string) bool {
	viewsDir := path.Join(baseDir, "views")
	files, err := fs.ReadDir(src, viewsDir)
	if err != nil {
		log.Error("读取模板目录失败：", err.Error())
		return false
	}

	registered := false
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".html") {
			continue
		}

		content, rErr := fs.ReadFile(src, path.Join(viewsDir, f.Name()))
		if rErr != nil {
			log.Error("读取模板文件失败：", rErr.Error())
			continue
		}

		name := f.Name()
		if namePrefix != "" {
			name = namePrefix + "/" + name
		}

		if _, pErr := tmpl.New(name).Parse(string(content)); pErr != nil {
			log.Error("解析模板失败：", pErr.Error())
			continue
		}

		registered = true
	}

	return registered
}

func registerAssetsFromFS(e *gin.Engine, src fs.FS, baseDir, routePath string) {
	assetsDir := path.Join(baseDir, "assets")
	if _, err := fs.Stat(src, assetsDir); err != nil {
		return
	}

	e.StaticFS(routePath, http.FS(subFS(src, assetsDir)))
}

func subFS(src fs.FS, dir string) fs.FS {
	sub, _ := fs.Sub(src, dir)

	return sub
}
