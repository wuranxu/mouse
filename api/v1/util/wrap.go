package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func Error(ctx *gin.Context, code int, msg error) {
	ctx.JSON(200, &MouseResponse{
		Code: code,
		Data: nil,
		Msg:  msg.Error(),
	})
}

func Wrap(f func(ctx *gin.Context) (any, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := f(c)
		code, ok := data.(int)
		if !ok {
			code = 110
		}
		if err != nil {
			Error(c, code, err)
			return
		}
		Success(c, data)
	}
}

func ErrWrap(err error, msg error) error {
	return fmt.Errorf("%v: %v", err, msg)
}
