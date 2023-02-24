package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api"
)

func main() {
	app := gin.Default()
	app.Use(gin.Logger(), gin.Recovery())

	// register route
	api.Register(app)

	app.Run(":8881")
}
