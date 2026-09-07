package api

import (
	"nimbus/api/handlers"
	"nimbus/client"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")

	newAuthenticator := func() handlers.Authenticator {
		return client.NewClient()
	}

	v1.POST("/login", handlers.LoginHandler(newAuthenticator))

	r.GET("/health", handlers.HealthHandler)

	return r
}
