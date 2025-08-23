package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/utils"
)

type RateLimiter struct {
    mu       sync.RWMutex
    clients  map[string]*Client
    requests int
    window   time.Duration
}

type Client struct {
    requests int
    window   time.Time
}

func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
    rl := &RateLimiter{
        clients:  make(map[string]*Client),
        requests: requests,
        window:   window,
    }
    
    // Cleanup goroutine
    go rl.cleanup()
    
    return rl
}

func (rl *RateLimiter) Allow(clientID string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    client, exists := rl.clients[clientID]
    now := time.Now()
    
    if !exists || now.Sub(client.window) > rl.window {
        rl.clients[clientID] = &Client{
            requests: 1,
            window:   now,
        }
        return true
    }
    
    if client.requests >= rl.requests {
        return false
    }
    
    client.requests++
    return true
}

func (rl *RateLimiter) cleanup() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            rl.mu.Lock()
            now := time.Now()
            for clientID, client := range rl.clients {
                if now.Sub(client.window) > rl.window {
                    delete(rl.clients, clientID)
                }
            }
            rl.mu.Unlock()
        }
    }
}

// Global rate limiter instance
var globalRateLimiter = NewRateLimiter(100, time.Minute) // 100 requests per minute

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientID := c.ClientIP()
        
        if !globalRateLimiter.Allow(clientID) {
            utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests")
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Custom rate limiter for specific endpoints
func CustomRateLimitMiddleware(requests int, window time.Duration) gin.HandlerFunc {
    limiter := NewRateLimiter(requests, window)
    
    return func(c *gin.Context) {
        clientID := c.ClientIP()
        
        if !limiter.Allow(clientID) {
            utils.ErrorResponse(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests")
            c.Abort()
            return
        }
        
        c.Next()
    }
}
