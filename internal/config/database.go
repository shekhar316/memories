package config

import (
    "github.com/shekhar316/memories/pkg/database"
    "github.com/shekhar316/memories/pkg/logger"
)

func (c *Config) InitializeDatabase() (*database.Database, error) {
    logger.Sugar.Info("Initializing database connection...")
    
    db, err := database.NewDatabase(
        c.Database.Host,
        c.Database.Port,
        c.Database.User,
        c.Database.Password,
        c.Database.DBName,
        c.Database.SSLMode,
        c.Server.Mode,
    )
    
    if err != nil {
        logger.Sugar.Errorw("Failed to initialize database", "error", err)
        return nil, err
    }

    // Run migrations
    if c.Database.AutoMigrate == true {
        if err := db.AutoMigrate(); err != nil {
            logger.Sugar.Errorw("Failed to run database migrations", "error", err)
            return nil, err
        }
    } else {
        logger.Sugar.Info("Skipping database migrations")
    }
    

    return db, nil
}
