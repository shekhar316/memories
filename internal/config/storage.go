package config

import (
    "context"
    
    "github.com/shekhar316/memories/pkg/storage"
    "github.com/shekhar316/memories/pkg/logger"
)

func (c *Config) InitializeS3Storage(ctx context.Context) (*storage.S3Storage, error) {
    logger.Sugar.Info("Initializing S3 storage...")

    s3Config := storage.S3Config{
        Region:          c.Storage.AWSRegion,
        Bucket:          c.Storage.AWSBucket,
        AccessKeyID:     c.Storage.AWSAccessKey,
        SecretAccessKey: c.Storage.AWSSecretKey,
        EndpointURL:     c.Storage.EndpointURL,
        ForcePathStyle:  c.Storage.ForcePathStyle,
    }

    s3Storage, err := storage.NewS3Storage(ctx, s3Config)
    if err != nil {
        logger.Sugar.Errorw("Failed to initialize S3 storage", "error", err)
        return nil, err
    }

    // Test connection immediately
    if err := s3Storage.TestConnection(ctx); err != nil {
        logger.Sugar.Errorw("S3 connection test failed", "error", err)
        return nil, err
    }

    logger.Sugar.Info("S3 storage initialized and connection tested successfully")
    return s3Storage, nil
}
