package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/utils"
    "github.com/shekhar316/memories/pkg/database"
)


type HealthHandler struct {
    db *database.Database
}

func NewHealthHandler(db *database.Database) *HealthHandler {
    return &HealthHandler{
        db: db,
    }
}

type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Service   string    `json:"service"`
    Version   string    `json:"version"`
    Uptime    string    `json:"uptime"`
    Database  map[string]interface{} `json:"database,omitempty"`
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


func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
    uptime := time.Since(startTime)
    status := "ready"
    statusCode := http.StatusOK
    
    dbStatus := make(map[string]interface{})
    
    // Check database connection
    if h.db != nil {
        if err := h.db.Ping(); err != nil {
            dbStatus["status"] = "disconnected"
            dbStatus["error"] = err.Error()
            status = "not ready"
            statusCode = http.StatusServiceUnavailable
        } else {
            dbStatus["status"] = "connected"
            dbStatus["stats"] = h.db.GetStats()
        }
    } else {
        dbStatus["status"] = "not initialized"
        status = "not ready"
        statusCode = http.StatusServiceUnavailable
    }

    healthData := HealthResponse{
        Status:    status,
        Timestamp: time.Now(),
        Service:   "memories-backend",
        Version:   "1.0.0",
        Uptime:    uptime.String(),
        Database:  dbStatus,
    }

    // TODO: Implement storage check

    
    utils.SuccessResponse(c, statusCode, "Service readiness check.", healthData)
}