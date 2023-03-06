package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/middleware"
	"github.com/wuranxu/mouse/service/auth"
	"github.com/wuranxu/mouse/utils/request"
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute(middlewares ...gin.HandlerFunc) {
	group := a.app.Group("/auth", middlewares...)

	group.POST("/login", middleware.ValidateJSON[dto.LoginDto], request.Wrap(auth.Login))
	group.POST("/register", middleware.ValidateJSON[dto.RegisterDto], request.Wrap(auth.Register))
	group.GET("/currentUser", request.Wrap(auth.Query))
}
