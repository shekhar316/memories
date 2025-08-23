package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/utils"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Service   string    `json:"service"`
    Version   string    `json:"version"`
    Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

func (h *HealthHandler) HealthCheck(c *gin.Context) {
    uptime := time.Since(startTime)
    
    healthData := HealthResponse{
        Status:    "healthy",
        Timestamp: time.Now(),
        Service:   "memories-backend",
        Version:   "1.0.0",
        Uptime:    uptime.String(),
    }
    
    utils.SuccessResponse(c, http.StatusOK, "Application is running fine.", healthData)
}
