package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Authenticate(username, password string) error
}

type AuthenticatorFactory func() Authenticator

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(newAuthenticator AuthenticatorFactory) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// A fresh client means a fresh cookie jar. Sharing one Librus client
		// between requests would mix sessions belonging to different users.
		authenticator := newAuthenticator()
		err := authenticator.Authenticate(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged in successfully",
		})
	}
}
