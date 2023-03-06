package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/utils/request"
)

const (
	TokenCheckFailedCode = iota + 40010
	TokenNotFoundCode
)

var (
	TokenCheckFailed = errors.New("token expired")
	TokenNotFound    = errors.New("token not found")
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		author, err := context.Cookie(conf.MouseToken)
		if err != nil || author == "" {
			request.Failed(context, TokenNotFoundCode, TokenNotFound)
			context.Abort()
			return
		}
		user, err := JWTUtil.ParseToken(author)
		if err != nil {
			request.Failed(context, TokenCheckFailedCode, TokenCheckFailed)
			context.Abort()
			return
		}
		context.Set(conf.MouseUser, user.MouseUser)
		context.Next()
	}
}
