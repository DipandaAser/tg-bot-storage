package handlers

import (
	"github.com/DipandaAser/tg-bot-storage/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

//ApiKeyMiddleware is a middleware that checks if the request has a valid api key
func ApiKeyMiddleware(c *gin.Context) {
	apiKey := c.Request.Header.Get("api-key")
	if apiKey != config.GetDefaultConfig().ApiKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad api key"})
		c.Abort()
	}
}
