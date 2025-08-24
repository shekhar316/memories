package database

import (
    "fmt"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/shekhar316/memories/internal/models"
    pkgLogger "github.com/shekhar316/memories/pkg/logger"
)

var DB *gorm.DB

type Database struct {
    DB *gorm.DB
}

func NewDatabase(host, port, user, password, dbname, sslmode, ginMode string) (*Database, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
        host, user, password, dbname, port, sslmode,
    )

    // Configure GORM logger
    var gormLogger logger.Interface
    if ginMode == "debug" {
        gormLogger = logger.Default.LogMode(logger.Info)
    } else {
        gormLogger = logger.Default.LogMode(logger.Silent)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })

    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get database instance: %w", err)
    }

    // Connection pool settings
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    sqlDB.SetConnMaxIdleTime(1 * time.Minute)

    // Set global DB instance
    DB = db

    pkgLogger.Sugar.Info("Database connection established successfully")

    return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
    pkgLogger.Sugar.Info("Starting database migration...")
    
    // Enable UUID extension
    if err := d.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
        return fmt.Errorf("failed to create uuid extension: %w", err)
    }

    // Auto migrate all models
    err := d.DB.AutoMigrate(
        &models.User{},
        &models.SubscriptionModel{},
        &models.UserSubscription{},
        &models.Album{},
        &models.AlbumShare{},
        &models.Photo{},
        &models.ProfilePhoto{},
        &models.Transaction{},
    )

    if err != nil {
        return fmt.Errorf("failed to auto migrate: %w", err)
    }

    // Create indexes
    if err := d.createCustomIndexes(); err != nil {
        return fmt.Errorf("failed to create custom indexes: %w", err)
    }

    pkgLogger.Sugar.Info("Database migration completed successfully")
    return nil
}

func (d *Database) createCustomIndexes() error {
    indexes := []string{
        "CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email) WHERE is_active = 1",
    }

    for _, index := range indexes {
        if err := d.DB.Exec(index).Error; err != nil {
            pkgLogger.Sugar.Warnw("Failed to create index", "index", index, "error", err)
        }
    }

    return nil
}


func (d *Database) Ping() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}

func (d *Database) Close() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (d *Database) GetStats() map[string]interface{} {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return map[string]interface{}{
            "error": err.Error(),
        }
    }

    stats := sqlDB.Stats()
    return map[string]interface{}{
        "max_open_connections": stats.MaxOpenConnections,
        "open_connections":     stats.OpenConnections,
        "in_use":              stats.InUse,
        "idle":                stats.Idle,
    }
}
