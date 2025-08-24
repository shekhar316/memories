package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/handlers"
    "github.com/shekhar316/memories/pkg/database"
    "github.com/shekhar316/memories/pkg/storage"
)

type RouteHandler struct {
    HealthHandler *handlers.HealthHandler
}

func NewRouteHandler(db *database.Database, s3Storage *storage.S3Storage) *RouteHandler {
    return &RouteHandler{
        HealthHandler: handlers.NewHealthHandler(db, s3Storage),
        // Initialize other handlers when they're created
    }
}

func SetupRoutes(router *gin.Engine, db *database.Database, s3Storage *storage.S3Storage) {
	setupGlobalMiddleware(router)

    routeHandler := NewRouteHandler(db, s3Storage)

    // Setup individual route groups
    setupHealthRoutes(router, routeHandler.HealthHandler)
}
