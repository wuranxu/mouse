package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuranxu/mouse/api"
	"github.com/wuranxu/mouse/conf"
	"github.com/wuranxu/mouse/dao"
	"github.com/wuranxu/mouse/middleware"
	"log"
)

var (
	serverHost = flag.String("host", "0.0.0.0", "mouse server host")
	serverPort = flag.Int("port", 9527, "mouse server port")
	configPath = flag.String("config", "./conf.yml", "mouse config filepath")
)

func main() {
	flag.Parse()
	if err := conf.Init(*configPath); err != nil {
		log.Fatal("init config error: ", err)
	}
	if err := dao.InitDatabase(); err != nil {
		log.Fatal("create/update table failed: ", err)
	}
	app := gin.New()
	app.Use(middleware.Cors())
	app.Use(gin.Logger(), gin.Recovery())

	// register route
	api.Register(app)

	app.Run(fmt.Sprintf("%s:%d", *serverHost, *serverPort))
}
