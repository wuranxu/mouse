package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/utils/request"
)

const (
	ParametersCheckErrorCode = 10001 + iota
)

func ValidateJSON[T any](ctx *gin.Context) {
	params, err := request.ValidateJSON[T](ctx)
	if err != nil {
		ctx.Abort()
		request.Failed(ctx, ParametersCheckErrorCode, err)
		return
	}
	ctx.Set(conf.MouseParameters, params)
	ctx.Next()
}

func ValidateQuery[T any](ctx *gin.Context) {
	params, err := request.ValidateQuery[T](ctx)
	if err != nil {
		ctx.Abort()
		request.Failed(ctx, ParametersCheckErrorCode, err)
		return
	}
	ctx.Set(conf.MouseParameters, params)
	ctx.Next()
}
