package handlers

import (
	"net/http"
	"nimbus/session"
	"time"

	"github.com/gin-gonic/gin"
)

func TimetableHandler(sessions *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("nimbus_session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		librusClient, ok := sessions.Get(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		dateText := c.Query("weekStart")
		weekStart, err := time.Parse("2006-01-02", dateText)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		days, err := librusClient.GetNimbusTimetable(c.Request.Context(), weekStart)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch timetable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"days": days})
	}
}
