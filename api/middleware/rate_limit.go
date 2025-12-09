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

/*
* Create a new IP rate limiter
* @param r rate.Limit
* @param b int
* @return *IPRateLimiter
 */
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

/*
* Get rate limiter for a specific IP
* @param ip string
* @return *rate.Limiter
 */
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

/*
* Middleware to rate limit requests
* @param r rate.Limit
* @param b int
* @return gin.HandlerFunc
 */
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
