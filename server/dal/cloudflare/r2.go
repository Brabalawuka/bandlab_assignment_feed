package cloudflare

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Config holds the configuration for R2
type Config struct {
    AccessKey  string
    SecretKey  string
    AccountID  string
    BucketName string
}

// R2Service defines the interface for R2 operations
type R2Service interface {
    GetClient() *s3.Client
    GetPresignClient() *s3.PresignClient
    GetBucketName() string
    PresignPutObject(context context.Context, input *s3.PutObjectInput, expiration time.Duration) (string, error)
}

// R2Client is a wrapper for the R2 session and client
type R2Client struct {
    client        *s3.Client
    presignClient *s3.PresignClient
    bucketName    string
}

var (
    once     sync.Once
    r2Client R2Service
)

// Initialize sets up the R2 session and client
func Initialize(cfg *Config) error {
    var initError error
    once.Do(func() {
        if err := validateConfig(cfg); err != nil {
            initError = err
            hlog.Error(initError)
            return
        }

        awsCfg, err := config.LoadDefaultConfig(context.TODO(),
            config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
            config.WithRegion("auto"),
        )
        if err != nil {
            initError = fmt.Errorf("[r2service] failed to load AWS config: %w", err)
            hlog.Error(initError)
            return
        }

        client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
            o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
        })

        r2Client = &R2Client{
            client:        client,
            presignClient: s3.NewPresignClient(client),
            bucketName:    cfg.BucketName,
        }
    })
    return initError
}

func GetR2Service() R2Service {
	return r2Client
}

// validateConfig checks if the provided configuration is valid
func validateConfig(cfg *Config) error {
    if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.AccountID == "" || cfg.BucketName == "" {
        return fmt.Errorf("[r2service] invalid configuration: all fields must be non-empty")
    }
    return nil
}

// GetClient returns the global R2 client
func (c *R2Client) GetClient() *s3.Client {
    if c == nil {
        hlog.Error("[r2service] R2 client is not initialized")
        return nil
    }
    return c.client
}

// GetPresignClient returns the global R2 presign client
func (c *R2Client) GetPresignClient() *s3.PresignClient {
    if c == nil {
        hlog.Error("[r2service] R2 presign client is not initialized")
        return nil
    }
    return c.presignClient
}

// GetBucketName returns the bucket name
func (c *R2Client) GetBucketName() string {
    if c == nil {
        hlog.Error("[r2service] R2 client is not initialized")
        return ""
    }
    return c.bucketName
}

// PresignPutObject generates a presigned URL for a PutObject request
func (c *R2Client) PresignPutObject(context context.Context, input *s3.PutObjectInput, expiration time.Duration) (string, error) {
    presignResult, err := c.presignClient.PresignPutObject(context, input, s3.WithPresignExpires(expiration))

    if err != nil {
        return "", fmt.Errorf("[r2service] couldn't get presigned URL for PutObject: %w", err)
    }

    return presignResult.URL, nil
}

