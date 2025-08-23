package config

import (
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Storage  StorageConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Port    string
    Mode    string
	Timeout time.Duration
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type JWTConfig struct {
    Secret string
    Expiry time.Duration
}

type StorageConfig struct {
    AWSRegion    string
    AWSBucket    string
    AWSAccessKey string
    AWSSecretKey string
}

type LoggerConfig struct {
    Level string
}

func LoadConfig() (*Config, error) {
    godotenv.Load()

    config := &Config{
        Server: ServerConfig{
            Port:    getEnv("PORT", "8080"),
            Mode:    getEnv("GIN_MODE", "release"),
			Timeout: getDurationEnv("SERVER_TIMEOUT", 60*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", ""),
        },
        JWT: JWTConfig{
            Secret: getEnv("JWT_SECRET", ""),
            Expiry: getDurationEnv("JWT_EXPIRY", 24*time.Hour),
        },
        Storage: StorageConfig{
            AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
            AWSBucket:    getEnv("AWS_BUCKET", "ai-photos-storage"),
            AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
            AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
        },
        Logger: LoggerConfig{
            Level: getEnv("LOG_LEVEL", "info"),
        },
    }

    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}
