package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/handlers"
    "github.com/shekhar316/memories/pkg/database"
)

type RouteHandler struct {
    HealthHandler *handlers.HealthHandler
}

func NewRouteHandler(db *database.Database) *RouteHandler {
    return &RouteHandler{
        HealthHandler: handlers.NewHealthHandler(db),
        // Initialize other handlers when they're created
    }
}

func SetupRoutes(router *gin.Engine, db *database.Database) {
	setupGlobalMiddleware(router)

    routeHandler := NewRouteHandler(db)

    // Setup individual route groups
    setupHealthRoutes(router, routeHandler.HealthHandler)
}
