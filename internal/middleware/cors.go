package middleware

import (
    "time"
    
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
