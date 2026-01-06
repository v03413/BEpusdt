package router

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/go-cache"
)

var engine *gin.Engine
var authRoute = make(map[string]bool)
var secureRoute = make(map[string]struct{})

func Handler() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine = gin.New()

	session := memstore.NewStore([]byte(conf.Secret))
	session.Options(sessions.Options{
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

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

		ctx.HTML(200, "index.html", gin.H{"title": conf.Desc, "url": conf.Github})
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
		if conf.Debug {
			ctx.Next()
			return
		}

		var route = fmt.Sprintf("%s.%s", ctx.Request.Method, ctx.Request.URL.Path)
		if _, ok := secureRoute[route]; ok {
			sess := sessions.Default(ctx)
			if secure, ok := sess.Get(conf.AdminSecureK).(bool); !ok || !secure {
				ctx.JSON(403, gin.H{"code": 403, "msg": "unauthorized access"})
				ctx.Abort()
				return
			}
		}

		var need, ok = authRoute[route]
		if !ok || !need {
			ctx.Next()
			return
		}

		token, ok := cache.Get(conf.AdminTokenK)
		if !ok {
			ctx.JSON(403, gin.H{"code": 403, "msg": "token expired, please login again"})
			ctx.Abort()
			return
		}

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(403, gin.H{"code": 403, "msg": "missing authorization token"})
			ctx.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(cast.ToString(token)), []byte(authHeader)) != 1 {
			ctx.JSON(403, gin.H{"code": 403, "msg": "invalid authorization token"})
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
	var route = fmt.Sprintf("POST.%s%s", router.BasePath(), relativePath)

	authRoute[route] = checkAuth
	secureRoute[route] = struct{}{}

	router.POST(relativePath, handlers...)
}

func GetRegister(router *gin.RouterGroup, relativePath string, checkAuth bool, handlers ...gin.HandlerFunc) {
	var route = fmt.Sprintf("GET.%s%s", router.BasePath(), relativePath)

	authRoute[route] = checkAuth
	secureRoute[route] = struct{}{}

	router.GET(relativePath, handlers...)
}
