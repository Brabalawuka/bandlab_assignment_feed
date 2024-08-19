package cloudflare

import (
	"bandlab_feed_server/config"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var presignAllowedTypes = map[string]struct{}{
    "image/jpeg": {},
    "image/png": {},
    "image/bmp":  {},
}

// GeneratePresignedURL generates a presigned URL for uploading an object
func GeneratePresignedURL(context context.Context, filename string, filesize int64, contentType string) (string, error) {
    
    if  _, ok := presignAllowedTypes[contentType]; !ok {
        return "", fmt.Errorf("[GeneratePresignedURL] invalid content type: %s", contentType)
    }

    input := &s3.PutObjectInput{
        Bucket:       aws.String(GetR2Service().GetBucketName()),
        Key:          aws.String(filename),
        ContentType:  aws.String(contentType),
        ContentLength: aws.Int64(filesize),
    }

    return GetR2Service().PresignPutObject(context, input, time.Duration(config.AppConfig.R2ImagePresignExpirationSec) * time.Second)

}