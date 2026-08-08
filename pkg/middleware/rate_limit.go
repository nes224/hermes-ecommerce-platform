package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(requestsPerSec float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerSec), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error" : "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}
		c.Next()
	}
}