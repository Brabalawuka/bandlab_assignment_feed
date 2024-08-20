package handler

import (
	"bandlab_feed_server/config"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service"
	"bandlab_feed_server/util/async"
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func CreatePost(ctx context.Context, req *dto.CreatePostReq) (*dto.CreatePostResp, error) {

	// Create post in post service
	post, err := service.GetPostService().CreatePost(ctx, req)
	if err != nil {
		return nil, err
	}

	// Resize image and update post in post service, we use a goroutine to handle this for MVP
	// TODO: Use an MQ to handle this as failure should be retried
	if service.GetPostService().HasImage(req) {
		async.Go(ctx, "ResizeImageTask", func(ctx context.Context) {
			uploadedPath, err := service.GetImageService().ResizeAndUploadImage(ctx, req.ImageFilePath)
			if err != nil {
				hlog.CtxErrorf(ctx, "Error resizing and uploading image: %v", err)
				return
			}
			err = service.GetPostService().UpdatePostStatusAndImagePath(ctx, post.Id, uploadedPath)
			if err != nil {
				hlog.CtxErrorf(ctx, "Error updating post status and image path: %v", err)
			}
		})
	}

	hlog.CtxDebugf(ctx, "Successfully created post: %v, resp: %v", req, post)
	return post, nil
}

// func GetAllPostsCore(ctx context.Context, req GetAllPostsRequest) (*GetAllPostsResponse, error) {
// 	// Implement logic to fetch posts
// 	// For example:
// 	posts := []Post{} // Fetch from database
// 	return &GetAllPostsResponse{
// 		Posts:      posts,
// 		TotalCount: len(posts),
// 		NextPage:   req.Page + 1,
// 	}, nil
// }

func GetPresignedURLCore(ctx context.Context, req *dto.GetPresignedURLRequest) (*dto.GetPresignedURLResponse, error) {

	// Get presigned URL from image service
	url, filePath, err := service.GetImageService().GetPresignedURL(ctx, req.FileName, req.FileSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "Error getting presigned URL: %v", err)
		return nil, err
	}

	hlog.CtxDebugf(ctx, "Successfully got presigned URL: %v, filePath: %v", url, filePath)
	return &dto.GetPresignedURLResponse{
		URL:        url,
		ImagePath:  filePath,
		ExpiresIn:  config.AppConfig.R2ImagePresignExpirationSec,
	}, nil
}
