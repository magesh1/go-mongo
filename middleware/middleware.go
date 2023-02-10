package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/juju/ratelimit.v1"
)

// creating channel for rate limit
// var Limit = make(chan struct{}, 5) // here 5 will be req we can send 5 req/sec
var limit = ratelimit.NewBucketWithRate(5, 5)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// select {
		// case Limit <- struct{}{}:
		// 	defer func() {
		// 		<-Limit
		// 	}()
		// 	c.Next()
		// default:
		// 	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		// 		"error": "too many request",
		// 	})
		// }

		if limit.TakeAvailable(1) == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"Error": "Too Many Request"})
			return
		}

		c.Next()

	}

}
