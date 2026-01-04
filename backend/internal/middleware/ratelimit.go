package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client   *redis.Client
	requests int
	duration time.Duration
}

func NewRateLimiter(client *redis.Client, requests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		client:   client,
		requests: requests,
		duration: duration,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rl.client == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		ctx := context.Background()

		current, err := rl.client.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if current == 1 {
			rl.client.Expire(ctx, key, rl.duration)
		}

		if current > int64(rl.requests) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests. Please try again later.",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
