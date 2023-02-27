package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api/v1/auth"
	"github.com/wuranxu/mouse/api/v1/job"
)

type Router interface {
	AddRoute(middleware ...gin.HandlerFunc)
}

func Register(app *gin.Engine) {
	auth.New(app).AddRoute()

	// register job route
	job.New(app).AddRoute()
}
