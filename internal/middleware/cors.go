package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CORS returns middleware that sets permissive CORS headers so the SPA frontend
// (served from a different origin) can call the API from the browser. The allowed
// origin defaults to "*" and can be locked down with the CORS_ALLOW_ORIGIN env var.
func CORS() gin.HandlerFunc {
	origin := os.Getenv("CORS_ALLOW_ORIGIN")
	if origin == "" {
		origin = "*"
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
