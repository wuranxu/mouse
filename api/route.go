package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api/v1/job"
)

type Router interface {
	AddRoute()
}

func Register(app *gin.Engine) {
	// register job route
	job.New(app).AddRoute()

}
