package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/config"
	"bandlab_feed_server/dal/cloudflare"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var presignAllowedTypes = map[string]string{
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".bmp":  "image/bmp",
}

// ImageService 定义了图像操作的接口
type ImageService interface {
	GetProcessedFileURLById(id string) string
	GetPresignedURL(ctx context.Context, filename string, filesize int64) (url string, filePath string, err error)
	ResizeAndUploadImage(ctx context.Context, imagePath string) (uploadedPath string, err error)
	RawImageExists(ctx context.Context, imagePath string) (bool, error) // New method
}

var (
	once     sync.Once
	imageSrv ImageService
)

func InitImageService() {
	once.Do(func() {
		r2Service := cloudflare.GetR2Service()
		if r2Service == nil {
			panic("R2 service is not initialized")
		}
		imageSrv = &ImageServiceImpl{
			r2Service: r2Service,
		}
	})
}

func GetImageService() ImageService {
	return imageSrv
}

// ImageServiceImpl is the implementation of ImageService
type ImageServiceImpl struct {
	r2Service cloudflare.R2Service
}

// GetProcessedFilePathById returns the processed file path by the image Id
func (s *ImageServiceImpl) GetProcessedFileURLById(id string) string {
	return fmt.Sprintf("%s%s.jpg", config.AppConfig.R2ProcessedImageURL, id)
}

// GetOriginalFilePathById returns the original file path by the image Id
func (s *ImageServiceImpl) GetOriginalFilePathByFileName(fileName string) string {
	return fmt.Sprintf("original/%s", fileName)
}

// GetPresignedURL generates a presigned URL for uploading an object
func (s *ImageServiceImpl) GetPresignedURL(ctx context.Context, filename string, filesize int64) (url string, filePath string, err error) {
	fileBase := filepath.Base(filename)
	fileExt := filepath.Ext(fileBase)
	contentType := ""
	if allowedType, ok := presignAllowedTypes[fileExt]; !ok {
		hlog.CtxErrorf(ctx, "[ImageServiceImpl] presign invalid content type: %s", fileExt)
		return "", "", errs.ErrInvalidContentType
	} else {
		contentType = allowedType
	}
	newImageId := primitive.NewObjectID().Hex() 
	newFilePath := s.GetOriginalFilePathByFileName(newImageId + fileExt)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(cloudflare.GetR2Service().GetBucketName()),
		Key:           aws.String(newFilePath),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(filesize),
	}

	url, err = cloudflare.GetR2Service().PresignPutObject(ctx, input, time.Duration(config.AppConfig.R2ImagePresignExpirationSec)*time.Second)
	if err != nil {
		return "", "", err
	}
	return url, newFilePath, nil
}

// ResizeAndUploadImage resizes the image and uploads it to R2
// TODO: for MVP, skip the resize and format part due to time limitation
func (s *ImageServiceImpl) ResizeAndUploadImage(ctx context.Context, imagePath string) (string, error) {

	// Create a buffer to store the file content
	var buf bytes.Buffer

	// Download the file from R2
	err := cloudflare.GetR2Service().DownloadFile(ctx, imagePath, &buf)
	if err != nil {
		hlog.CtxErrorf(ctx, "[ImageServiceImpl] failed to download image from R2: %v", err)
		return "", errs.ErrR2ImageDownLoadFailed
	}

	// Resize and reformat the image
	// TODO: for MVP, skip the resize and format part due to time limitation

	// Prepare the upload path
	uploadPath := fmt.Sprintf("600x600/%s%s", filepath.Base(imagePath), filepath.Ext(imagePath))

	// Upload the file back to R2
	err = cloudflare.GetR2Service().UploadFile(ctx, &buf, uploadPath)
	if err != nil {
		hlog.CtxErrorf(ctx, "[ImageServiceImpl] failed to upload image to R2: ", err)
		return "", errs.ErrR2ImageUploadFailed
	}

	return uploadPath, nil
}

// ImageExists checks if an image exists in R2
func (s *ImageServiceImpl) RawImageExists(ctx context.Context, imagePath string) (bool, error) {
	// Check if the file exists
	exists, err := s.r2Service.FileExists(ctx, imagePath)
	if err != nil {
		hlog.CtxErrorf(ctx, "[ImageServiceImpl] failed to check if image exists in R2: ", err)
		return false, errs.ErrImageExistsCheckFailed
	}
	return exists, nil
}
