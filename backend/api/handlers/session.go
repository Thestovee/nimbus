package handlers

import (
	"net/http"
	"nimbus/session"

	"github.com/gin-gonic/gin"
)

func SessionHandler(sessions *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("nimbus_session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
			return
		}

		_, ok := sessions.Get(token)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"authenticated": true})
	}
}
