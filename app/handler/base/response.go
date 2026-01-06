package base

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/log"
)

type Result struct {
	Code  int         `json:"code,omitempty"`  // 状态码
	Msg   string      `json:"msg,omitempty"`   // 状态消息
	Total int64       `json:"total,omitempty"` // 总条数
	Data  interface{} `json:"data,omitempty"`  // 数据
}

func Typer(i any) string {
	var t string
	switch i.(type) {
	case string:
		t = "string"
	case int:
		t = "int"
	case int8:
		t = "int8"
	case int16:
		t = "int16"
	case int32:
		t = "int32"
	case int64:
		t = "int64"
	case uint:
		t = "uint"
	case uint8:
		t = "uint8"
	case uint16:
		t = "uint16"
	case uint32:
		t = "uint32"
	case uint64:
		t = "uint64"
	case map[string]string:
		t = "map[string]string"
	case map[int]string:
		t = "map[int]string"
	case map[string]int:
		t = "map[string]int"
	case map[int]int:
		t = "map[int]int"
	case []string:
		t = "[]string"
	case []int:
		t = "[]int"
	case [1]string:
		t = "array [1] string"
	case [2]string:
		t = "array [1] string"
	case [1]int:
		t = "array [1] int"
	case [2]int:
		t = "array [2] int"
	default:
		t = "unknow"
	}
	return t
}

func Response(ctx *gin.Context, code int, data ...any) {
	var d = Result{Code: code}

	if len(data) > 0 {
		if Typer(data[0]) == "string" {
			d.Msg = data[0].(string)
		} else {
			d.Data = data[0]
		}
		if len(data) == 2 {
			d.Total = data[1].(int64)
		}
	}

	if code >= 500 && code <= 599 {
		if len(data) > 0 {
			err, ok := data[0].(error)
			if ok {
				log.Error(err.Error())
			} else {
				log.Error(fmt.Sprintf("%v", data[0]))
			}
		}

		d.Msg = "服务器错误"
	}

	ctx.JSON(200, d)
}

func Ok(ctx *gin.Context, data ...any) {
	Response(ctx, 200, data...)
}

// Error 500 服务器错误
func Error(ctx *gin.Context, err error) {
	Response(ctx, 500, err)
}

// BadRequest 400 请求错误
func BadRequest(ctx *gin.Context, err string) {
	Response(ctx, 400, err)
}
