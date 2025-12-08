package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	i := NewIPRateLimiter(r, b)
	return func(c *gin.Context) {
		limiter := i.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  false,
				"message": "Too many requests, slow down!",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
