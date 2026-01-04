package router

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
)

var engine *gin.Engine
var authRoute = make(map[string]bool)

func Handler() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine = gin.New()

	session := memstore.NewStore([]byte(conf.Secret))
	session.Options(sessions.Options{MaxAge: 86400, HttpOnly: true})

	engine.Use(sessions.Sessions("session", session))
	engine.Use(gin.LoggerWithWriter(log.GetWriter()), gin.Recovery())
	engine.Use(sessionAuth(), copyright())
	engine.NoRoute(noRoute())
	engine.GET("/", func(ctx *gin.Context) {
		sess := sessions.Default(ctx)
		if secure, ok := sess.Get(conf.AdminSecureK).(bool); ok && secure {
			ctx.HTML(200, "secure.html", gin.H{})

			return
		}

		ctx.HTML(200, "index.html", gin.H{"title": "一款更好用的加密货币收款网关", "url": conf.Github})
	})

	{
		staticInit(engine)
		epusdtInit(engine)
		epayInit(engine)
		adminInit(engine)
		authInit(engine)
	}

	return engine
}

func sessionAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var need, ok = authRoute[ctx.Request.Method+"."+ctx.Request.URL.Path]
		if !ok || !need {
			ctx.Next() // 不需要鉴权

			return
		}

		session := sessions.Default(ctx)
		token, ok := session.Get(conf.AdminTokenK).(string)
		if !ok || token == "" || token != ctx.GetHeader("Authorization") {
			//开发的时候可以注释掉，方便调试
			ctx.JSON(403, gin.H{"code": 403, "msg": "invalid token"})
			ctx.Abort()

			return
		}

		ctx.Next()
	}
}

func noRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == model.GetC(model.AdminSecure) {
			session := sessions.Default(ctx)
			session.Set(conf.AdminSecureK, true)
			_ = session.Save()

			ctx.Redirect(302, "/#/login")

			return
		}
	}
}

func copyright() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Payment-Gateway", "https://github.com/v03413/BEpusdt")
	}
}

func PostRegister(router *gin.RouterGroup, relativePath string, checkAuth bool, handlers ...gin.HandlerFunc) {
	var path = router.BasePath() + relativePath

	authRoute[fmt.Sprintf("POST.%s", path)] = checkAuth

	router.POST(relativePath, handlers...)
}

func GetRegister(router *gin.RouterGroup, relativePath string, checkAuth bool, handlers ...gin.HandlerFunc) {
	var path = router.BasePath() + relativePath

	authRoute[fmt.Sprintf("GET.%s", path)] = checkAuth

	router.GET(relativePath, handlers...)
}
