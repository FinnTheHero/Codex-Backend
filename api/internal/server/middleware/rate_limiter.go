package firestore_middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 30)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests!",
			})
			return
		}
		// TODO: Improve rate limiting logic. Add more detailed cooldown to IPs itself.
		c.Next()
	}
}
