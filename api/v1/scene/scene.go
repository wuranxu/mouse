package scene

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/dto"
	"github.com/wuranxu/mouse/middleware"
	"github.com/wuranxu/mouse/service/scene"
	"github.com/wuranxu/mouse/utils/request"
)

type Api struct {
	app *gin.Engine
}

func New(app *gin.Engine) *Api {
	return &Api{app: app}
}

func (a *Api) AddRoute(middlewares ...gin.HandlerFunc) {
	group := a.app.Group("/scene", middlewares...)
	group.GET("/", request.Handle(scene.QueryScene))
	group.GET("/list", request.Handle(scene.ListScene))
	group.POST("/insert", middleware.ValidateJSON[dto.SceneDto], request.Handle(scene.CreateScene))
}
