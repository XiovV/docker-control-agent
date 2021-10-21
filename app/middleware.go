package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("key")

		if !app.config.CompareHash(apiKey) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid api key"})
			return
		}

		c.Next()
	}
}
