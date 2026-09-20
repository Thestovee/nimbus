package handlers

import (
	"nimbus/client"
	"nimbus/session"

	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginFunc func(username, password string) (*client.LibrusClient, error)
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(login LoginFunc, sessions *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// A fresh client means a fresh cookie jar. Sharing one Librus client
		// between requests would mix sessions belonging to different users.
		librusClient, err := login(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		token, err := sessions.Create(librusClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "nimbus_session",
			Value:    token,
			Path:     "/",
			MaxAge:   24 * 60 * 60,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   false, // only for local testing, TODO don't forget to turn on when finishing
		})

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in successfully",
		})
	}
}
