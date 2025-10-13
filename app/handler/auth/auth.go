package auth

import (
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/handler/base"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/utils"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
}

type authLoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type authPasswordReq struct {
	Password        string `json:"password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

var types = map[string]string{
	model.TradeTypeTronTrx:      "TRON・TRX",
	model.TradeTypeUsdtTrc20:    "USDT・TRC20",
	model.TradeTypeUsdtErc20:    "USDT・ERC20",
	model.TradeTypeUsdtBep20:    "USDT・BEP20",
	model.TradeTypeUsdtAptos:    "USDT・APTOS",
	model.TradeTypeUsdtXlayer:   "USDT・Xlayer",
	model.TradeTypeUsdtSolana:   "USDT・Solana",
	model.TradeTypeUsdtPolygon:  "USDT・Polygon",
	model.TradeTypeUsdtArbitrum: "USDT・Arbitrum",
	model.TradeTypeUsdcErc20:    "USDC・ERC20",
	model.TradeTypeUsdcBep20:    "USDC・BEP20",
	model.TradeTypeUsdcXlayer:   "USDC・Xlayer",
	model.TradeTypeUsdcPolygon:  "USDC・Polygon",
	model.TradeTypeUsdcArbitrum: "USDC・Arbitrum",
	model.TradeTypeUsdcBase:     "USDC・Base",
	model.TradeTypeUsdcTrc20:    "USDC・TRC20",
	model.TradeTypeUsdcSolana:   "USDC・Solana",
	model.TradeTypeUsdcAptos:    "USDC・APTOS",
}

func (Auth) Info(ctx *gin.Context) {
	base.Ok(ctx, gin.H{
		"admin_username": model.GetK(model.AdminUsername),
		"trade_type":     types,
		"trade_fiat":     model.SupportTradeFiat,
		"trade_crypto":   model.SupportTradeCrypto,
	})
}

func (Auth) Menu(ctx *gin.Context) {
	type meta struct {
		Title     string   `json:"title"`
		Hide      bool     `json:"hide"`
		Disable   bool     `json:"disable"`
		KeepAlive bool     `json:"keepAlive"`
		Affix     bool     `json:"affix"`
		Link      string   `json:"link"`
		Iframe    bool     `json:"iframe"`
		IsFull    bool     `json:"isFull"`
		Roles     []string `json:"roles"`
		SvgIcon   string   `json:"svgIcon"`
		Icon      string   `json:"icon"`
		Sort      int      `json:"sort"`
		Type      int      `json:"type"`
	}
	type menu struct {
		Id        string `json:"id"`
		ParentId  string `json:"parentId"`
		Path      string `json:"path"`
		Name      string `json:"name"`
		Component string `json:"component"`
		Meta      meta   `json:"meta"`
		Children  []menu `json:"children"`
	}

	var data = []menu{
		{
			Id:        "01",
			Path:      "/home",
			Name:      "home",
			Component: "home/home",
			Meta: meta{
				Title:     "home",
				Hide:      false,
				Disable:   false,
				KeepAlive: false,
				Affix:     true,
				Link:      "",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "home",
				Icon:      "",
				Sort:      1,
				Type:      2,
			},
			Children: nil,
		},
		{
			Id:        "02",
			Path:      "/wallet",
			Name:      "wallet",
			Component: "wallet/wallet",
			Meta: meta{
				Title:     "wallet",
				Hide:      false,
				Disable:   false,
				KeepAlive: true,
				Affix:     false,
				Link:      "",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "classify",
				Icon:      "",
				Sort:      1,
				Type:      2,
			},
			Children: nil,
		},
		{
			Id:        "03",
			Path:      "/order",
			Name:      "order",
			Component: "order/order",
			Meta: meta{
				Title:     "order",
				Hide:      false,
				Disable:   false,
				KeepAlive: true,
				Affix:     false,
				Link:      "",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "table",
				Icon:      "",
				Sort:      1,
				Type:      2,
			},
			Children: nil,
		},
		{
			Id:        "04",
			Path:      "/rate",
			Name:      "rate",
			Component: "rate/rate",
			Meta: meta{
				Title:     "rate",
				Hide:      false,
				Disable:   false,
				KeepAlive: true,
				Affix:     false,
				Link:      "",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "directives",
				Icon:      "",
				Sort:      1,
				Type:      1,
			},
			Children: []menu{
				{
					Id:        "0401",
					ParentId:  "04",
					Path:      "/rate/list",
					Name:      "rate-list",
					Component: "rate/list",
					Meta: meta{
						Title:     "rate-list",
						Hide:      false,
						Disable:   false,
						KeepAlive: true,
						Affix:     false,
						Link:      "",
						Iframe:    false,
						IsFull:    false,
						Roles:     []string{"admin"},
						SvgIcon:   "",
						Icon:      "icon-list",
						Sort:      1,
						Type:      2,
					},
					Children: nil,
				},
				{
					Id:        "0402",
					ParentId:  "04",
					Path:      "/rate/syntax",
					Name:      "rate-syntax",
					Component: "rate/syntax",
					Meta: meta{
						Title:     "rate-syntax",
						Hide:      false,
						Disable:   false,
						KeepAlive: true,
						Affix:     false,
						Link:      "",
						Iframe:    false,
						IsFull:    false,
						Roles:     []string{"admin"},
						SvgIcon:   "",
						Icon:      "icon-settings",
						Sort:      1,
						Type:      2,
					},
					Children: nil,
				},
			},
		},
		{
			Id:        "05",
			Path:      "/system",
			Name:      "system",
			Component: "system/system",
			Meta: meta{
				Title:     "system",
				Hide:      false,
				Disable:   false,
				KeepAlive: true,
				Affix:     false,
				Link:      "",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "set",
				Icon:      "",
				Sort:      1,
				Type:      1,
			},
			Children: []menu{
				{
					Id:        "0501",
					ParentId:  "05",
					Path:      "/system/base/base",
					Name:      "system-base",
					Component: "system/base/base",
					Meta: meta{
						Title:     "system-base",
						Hide:      false,
						Disable:   false,
						KeepAlive: true,
						Affix:     false,
						Link:      "",
						Iframe:    false,
						IsFull:    false,
						Roles:     []string{"admin"},
						SvgIcon:   "",
						Icon:      "icon-layers",
						Sort:      1,
						Type:      2,
					},
					Children: nil,
				},
				{
					Id:        "0502",
					ParentId:  "05",
					Path:      "/system/rpc/rpc",
					Name:      "system-rpc",
					Component: "system/rpc/rpc",
					Meta: meta{
						Title:     "system-rpc",
						Hide:      false,
						Disable:   false,
						KeepAlive: true,
						Affix:     false,
						Link:      "",
						Iframe:    false,
						IsFull:    false,
						Roles:     []string{"admin"},
						SvgIcon:   "",
						Icon:      "icon-thunderbolt",
						Sort:      1,
						Type:      2,
					},
					Children: nil,
				},
			},
		},
		{
			Id:   "06",
			Path: "/github-api-doc",
			Name: "github-api-doc",
			Meta: meta{
				Title:     "github-api-doc",
				Hide:      false,
				Disable:   false,
				KeepAlive: true,
				Affix:     false,
				Link:      "https://github.com/v03413/BEpusdt/blob/main/docs/api.md",
				Iframe:    false,
				IsFull:    false,
				Roles:     []string{"admin"},
				SvgIcon:   "about",
				Icon:      "",
				Sort:      1,
				Type:      2,
			},
			Children: nil,
		},
	}

	base.Ok(ctx, data)
}

func (Auth) Login(ctx *gin.Context) {
	var req authLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.Response(ctx, 400, err.Error())

		return
	}

	var username = model.GetK(model.AdminUsername)
	if req.Username != username {
		base.Response(ctx, 400, "用户名或密码错误")

		return
	}

	var password = model.GetK(model.AdminPassword)
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)) != nil {
		base.Response(ctx, 400, "用户名或密码错误")

		return
	}

	rand, _ := utils.GenerateTradeId()

	var token = utils.StrSha256(rand + ctx.ClientIP())

	session := sessions.Default(ctx)
	session.Set(conf.AdminTokenK, token)
	defer session.Save()

	model.SetK(model.AdminLoginIP, ctx.ClientIP())
	model.SetK(model.AdminLoginAt, cast.ToString(time.Now().Format(time.DateTime)))

	base.Response(ctx, 200, gin.H{"token": token, "types": types})
}

func (Auth) Logout(ctx *gin.Context) {
	model.SetK(model.AdminToken, "")

	base.Response(ctx, 200, "退出成功")
}

func (Auth) SetPassword(ctx *gin.Context) {
	var req authPasswordReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, err.Error())

		return
	}

	var password = model.GetK(model.AdminPassword)
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)) != nil {
		base.BadRequest(ctx, "原密码错误")

		return
	}

	if req.ConfirmPassword != req.NewPassword {
		base.BadRequest(ctx, "两次输入的新密码不一致")

		return
	}

	var newPassword = strings.TrimSpace(req.NewPassword)
	if len(newPassword) < 6 {
		base.BadRequest(ctx, "新密码长度不能少于6位")

		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	model.SetK(model.AdminPassword, string(hash))
	model.SetK(model.AdminToken, "")

	base.Ok(ctx, "修改成功，请重新登录")
}
