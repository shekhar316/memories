package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/handlers"
)

type RouteHandler struct {
    HealthHandler *handlers.HealthHandler
}

func NewRouteHandler() *RouteHandler {
    return &RouteHandler{
        HealthHandler: handlers.NewHealthHandler(),
    }
}

func SetupRoutes(router *gin.Engine) {
	setupGlobalMiddleware(router)

    routeHandler := NewRouteHandler()

    // Setup individual route groups
    setupHealthRoutes(router, routeHandler.HealthHandler)
}
