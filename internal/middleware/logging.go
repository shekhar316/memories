package middleware

import (
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/pkg/logger"
)

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Calculate latency
        latency := time.Since(start)

        // Get status code
        statusCode := c.Writer.Status()

        // Log request details
        if raw != "" {
            path = path + "?" + raw
        }

        logger.Sugar.Infow("HTTP Request",
            "status", statusCode,
            "method", c.Request.Method,
            "path", path,
            "ip", c.ClientIP(),
            "user_agent", c.Request.UserAgent(),
            "latency", latency,
            "response_size", c.Writer.Size(),
        )
    }
}
