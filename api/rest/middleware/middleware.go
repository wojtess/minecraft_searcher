package middleware

import (
	"minecraft_searcher/api/rest/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(c *gin.Context) {
	err := token.TokenValid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
