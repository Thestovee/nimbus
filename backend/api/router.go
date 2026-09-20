package api

import (
	"nimbus/api/handlers"
	"nimbus/client"
	"nimbus/session"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	sessions := session.NewStore()
	v1 := r.Group("/api/v1")

	login := func(username, password string) (*client.LibrusClient, error) {
		librusClient := client.NewClient()
		if err := librusClient.Authenticate(username, password); err != nil {
			return nil, err
		}
		return librusClient, nil
	}

	v1.POST("/login", handlers.LoginHandler(login, sessions))
	v1.GET("/session", handlers.SessionHandler(sessions))
	v1.GET("/timetable", handlers.TimetableHandler(sessions))

	r.GET("/health", handlers.HealthHandler)

	return r
}
