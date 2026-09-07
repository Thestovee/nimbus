package api

import (
	"nimbus/api/handlers"
	"nimbus/client"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", handlers.LoginHandler(client.NewClient))
	}

	r.GET("/health", handlers.HealthHandler)

	return r
}
