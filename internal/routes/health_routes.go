package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/handlers"
)

func setupHealthRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler) {
    
    // add to v1 API group for consistency
    v1 := router.Group("/api/v1")
    {
        v1.GET("/health", healthHandler.HealthCheck)
        v1.GET("/readiness", healthHandler.ReadinessCheck)
    }
}
