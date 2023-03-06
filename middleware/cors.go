package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"OPTION", "GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders: []string{"*"},
	})
}
