package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/utils"
    "github.com/shekhar316/memories/pkg/database"
    "github.com/shekhar316/memories/pkg/storage"
)


type HealthHandler struct {
    db *database.Database
    s3Storage *storage.S3Storage
}

func NewHealthHandler(db *database.Database, storage *storage.S3Storage) *HealthHandler {
    return &HealthHandler{
        db: db,
        s3Storage: storage,
    }
}

type HealthResponse struct {
    Status    string    `json:"status"`
    Timestamp time.Time `json:"timestamp"`
    Service   string    `json:"service"`
    Version   string    `json:"version"`
    Uptime    string    `json:"uptime"`
    Database  map[string]interface{} `json:"database,omitempty"`
    Storage   map[string]interface{} `json:"storage,omitempty"`
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
    storageStatus := make(map[string]interface{})
    
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


    if h.s3Storage != nil {
        if err := h.s3Storage.TestConnection(c.Request.Context()); err != nil {
            storageStatus["status"] = "disconnected"
            storageStatus["error"] = err.Error()
            status = "not ready"
            statusCode = http.StatusServiceUnavailable
        } else {
            storageStatus["status"] = "connected"
            storageStatus["stats"] = h.s3Storage.GetStats(c.Request.Context())
        }
    } else {
        storageStatus["status"] = "not initialized"
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
        Storage:   storageStatus,
    }
    
    utils.SuccessResponse(c, statusCode, "Service readiness check.", healthData)
}