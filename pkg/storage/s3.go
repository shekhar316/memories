package storage

import (
    "context"
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    pkgLogger "github.com/shekhar316/memories/pkg/logger"
)

type S3Storage struct {
    Client *s3.Client
    Bucket string
    Region string
}

type S3Config struct {
    Region          string
    Bucket          string
    AccessKeyID     string
    SecretAccessKey string
    EndpointURL     string
    ForcePathStyle  bool
}

func NewS3Storage(ctx context.Context, s3Config S3Config) (*S3Storage, error) {
    if s3Config.Bucket == "" {
        return nil, fmt.Errorf("S3 bucket name is required")
    }

    if s3Config.AccessKeyID == "" || s3Config.SecretAccessKey == "" {
        return nil, fmt.Errorf("AWS credentials are required")
    }

    // Configure AWS SDK
    configOptions := []func(*config.LoadOptions) error{
        config.WithRegion(s3Config.Region),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                s3Config.AccessKeyID,
                s3Config.SecretAccessKey,
                "",
            ),
        ),
    }

    awsCfg, err := config.LoadDefaultConfig(ctx, configOptions...)
    if err != nil {
        return nil, fmt.Errorf("failed to load AWS config: %w", err)
    }

    // Create S3 client with optional endpoint URL (for MinIO/LocalStack)
    clientOptions := []func(*s3.Options){
        func(o *s3.Options) {
            if s3Config.EndpointURL != "" {
                o.BaseEndpoint = aws.String(s3Config.EndpointURL)
            }
            if s3Config.ForcePathStyle {
                o.UsePathStyle = true
            }
        },
    }

    client := s3.NewFromConfig(awsCfg, clientOptions...)

    storage := &S3Storage{
        Client: client,
        Bucket: s3Config.Bucket,
        Region: s3Config.Region,
    }

    pkgLogger.Sugar.Infow("S3 storage client initialized", 
        "bucket", s3Config.Bucket, 
        "region", s3Config.Region)

    return storage, nil
}

// TestConnection tests S3 connection
func (s *S3Storage) TestConnection(ctx context.Context) error {
    // Create a context with timeout for the health check
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    _, err := s.Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
        Bucket: aws.String(s.Bucket),
    })

    if err != nil {
        // This is a very lightweight operation that tests both connectivity and permissions
        _, listErr := s.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
            Bucket:  aws.String(s.Bucket),
            MaxKeys: aws.Int32(1),
        })

        if listErr != nil {
            return fmt.Errorf("S3 connection test failed: %w", listErr)
        }
    }
    
    return nil
}

// GetBucketInfo returns basic bucket information for health checks
func (s *S3Storage) GetBucketInfo(ctx context.Context) (map[string]interface{}, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    info := map[string]interface{}{
        "bucket": s.Bucket,
        "region": s.Region,
    }

    // Try to get bucket location
    if location, err := s.Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
        Bucket: aws.String(s.Bucket),
    }); err == nil {
        if location.LocationConstraint != "" {
            info["actual_region"] = string(location.LocationConstraint)
        } else {
            info["actual_region"] = "us-east-1" // Default region
        }
    }

    // Get bucket versioning status (optional)
    if versioning, err := s.Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
        Bucket: aws.String(s.Bucket),
    }); err == nil {
        info["versioning"] = string(versioning.Status)
    }

    return info, nil
}

// Ping is an alias for TestConnection for consistency with database interface
func (s *S3Storage) Ping(ctx context.Context) error {
    return s.TestConnection(ctx)
}

// Close gracefully closes the S3 client (no-op for AWS SDK v2)
func (s *S3Storage) Close() error {
    // AWS SDK v2 doesn't require explicit closing
    pkgLogger.Sugar.Info("S3 storage client closed")
    return nil
}

// GetStats returns connection statistics
func (s *S3Storage) GetStats(ctx context.Context) map[string]interface{} {
    stats := map[string]interface{}{
        "bucket":     s.Bucket,
        "region":     s.Region,
        "client_type": "aws-sdk-go-v2",
    }

    // Try to get additional bucket info
    if info, err := s.GetBucketInfo(ctx); err == nil {
        for k, v := range info {
            stats[k] = v
        }
    } else {
        stats["info_error"] = err.Error()
    }

    return stats
}

// IsAccessible checks if the bucket is accessible
func (s *S3Storage) IsAccessible(ctx context.Context) bool {
    return s.TestConnection(ctx) == nil
}
