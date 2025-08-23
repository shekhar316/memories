package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/middleware"
)

func setupGlobalMiddleware(router *gin.Engine) {
    // Recovery middleware
    router.Use(gin.Recovery())
    
    // Logging middleware
    router.Use(middleware.LoggingMiddleware())
    
    // CORS middleware
    router.Use(middleware.CORSMiddleware())
    
    // Rate limiting middleware 
    router.Use(middleware.RateLimitMiddleware())
}

func setupProtectedMiddleware() []gin.HandlerFunc {
    return []gin.HandlerFunc{
        // middleware.AuthMiddleware(),
        // middleware.UserMiddleware(),
    }
}
