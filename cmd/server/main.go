package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/shekhar316/memories/internal/config"
	"github.com/shekhar316/memories/internal/routes"
    "github.com/shekhar316/memories/pkg/logger"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        panic(fmt.Sprintf("Failed to load config: %v", err))
    }

    // Initialize logger
    if err := logger.InitLogger(cfg.Logger.Level); err != nil {
        panic(fmt.Sprintf("Failed to initialize logger: %v", err))
    }

    defer logger.Sync()

    logger.Sugar.Info("Starting Memories Service...")


    ctx := context.Background()
    // Initialize database
    db, err := cfg.InitializeDatabase()
    if err != nil {
        logger.Sugar.Fatalf("Failed to initialize database: %v", err)
    }
    defer func() {
        if err := db.Close(); err != nil {
            logger.Sugar.Errorw("Failed to close database connection", "error", err)
        }
    }()

    s3Storage, err := cfg.InitializeS3Storage(ctx)
    if err != nil {
        logger.Sugar.Fatalf("Failed to initialize S3 storage: %v", err)
    }
    defer func() {
        if err := s3Storage.Close(); err != nil {
            logger.Sugar.Errorw("Failed to close S3 storage", "error", err)
        }
    }()

    // Set Gin mode
    gin.SetMode(cfg.Server.Mode)

    // Initialize Gin router
    router := gin.New()
	routes.SetupRoutes(router, db, s3Storage)

    // Create HTTP server
    server := &http.Server{
        Addr:         ":" + cfg.Server.Port,
        Handler:      router,
        ReadTimeout:  cfg.Server.Timeout,
        WriteTimeout: cfg.Server.Timeout,
    }

    // Start server in a goroutine
    go func() {
        logger.Sugar.Infof("Server starting on port %s", cfg.Server.Port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Sugar.Fatalf("Server failed to start: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Sugar.Info("Shutting down server...")

    // Give outstanding requests 30 seconds to complete
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        logger.Sugar.Fatalf("Server forced to shutdown: %v", err)
    }

    logger.Sugar.Info("Server exited")
}
