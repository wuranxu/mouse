package request

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/exception"
	"github.com/wuranxu/mouse/model"
	"github.com/wuranxu/mouse/utils/validate"
)

var (
	CheckParamsError    = exception.Err("参数错误")
	ParamsNotValidError = exception.Err("参数类型有误, 请检查")
)

const ok = "operate success"

type MouseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(200, &MouseResponse{
		Code: 0,
		Msg:  ok,
		Data: data,
	})
}

func Failed(ctx *gin.Context, code int, msg error) {
	ctx.JSON(200, &MouseResponse{
		Code: code,
		Data: nil,
		Msg:  msg.Error(),
	})
}

func Handle(f func(ctx *gin.Context) (any, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := f(c)
		if err != nil {
			code, ok := data.(int)
			if !ok {
				code = 110
			}
			Failed(c, code, err)
			return
		}
		Success(c, data)
	}
}

func ValidateJSON[T any](ctx *gin.Context) (*T, error) {
	var t T
	if err := ctx.ShouldBindJSON(&t); err != nil {
		return nil, err
	}
	if err := validate.Check(t, CheckParamsError); err != nil {
		return nil, err
	}
	return &t, nil
}

func ValidateQuery[T any](ctx *gin.Context) (*T, error) {
	var t T
	if err := ctx.ShouldBindQuery(&t); err != nil {
		return nil, ParamsNotValidError
	}
	if err := validate.Check(t, CheckParamsError); err != nil {
		return nil, err
	}
	return &t, nil
}

func GetUser(ctx *gin.Context) model.MouseUser {
	auth := ctx.MustGet(conf.MouseUser)
	return auth.(model.MouseUser)
}

func GetJson[T any](ctx *gin.Context) *T {
	data := ctx.MustGet(conf.MouseParameters)
	return data.(*T)
}
