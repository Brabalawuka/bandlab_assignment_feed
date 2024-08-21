package handler

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service"
	"bandlab_feed_server/util/async"
	"context"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleCreatePost(ctx context.Context, req *dto.CreatePostReq) (*dto.CreatePostResp, error) {

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


func HandleGetPresignedURL(ctx context.Context, req *dto.GetPresignedURLRequest) (*dto.GetPresignedURLResponse, error) {

	// Get presigned URL from image service
	resp, err := service.GetImageService().GetPresignedURL(ctx, req.FileName, req.FileSize)
	if err != nil {
		hlog.CtxErrorf(ctx, "Error getting presigned URL: %v", err)
		return nil, err
	}

	hlog.CtxDebugf(ctx, "Successfully got presigned URL: %v, filePath: %v", resp.URL, resp.ImagePath)
	return resp, nil
}

func HandleGetPost(ctx context.Context, req *dto.FetchPostsReq) (*dto.FetchPostsResp, error) {
	var posts []*dao.Post
	var hasMore bool
	var nextCursor string
	switch req.OrderBy {
	case dto.OrderByPostID:
		var err error
		var postID *primitive.ObjectID
		if req.PreviousCursor != "" {
			objectId, err := primitive.ObjectIDFromHex(req.PreviousCursor)
			if err != nil {
				hlog.CtxErrorf(ctx, "[HandleGetPost] invalid previous cursor: %v, when fetching by post id", err)
				return nil, errs.ErrInvalidRequest
			}
			postID = &objectId
		}
		posts, hasMore, err = service.GetPostService().FetchPostsByPostIDCursor(ctx, req.Limit, postID)
		if err != nil {
			hlog.CtxErrorf(ctx, "[HandleGetPost] error fetching posts: %v", err)
			return nil, err
		}
		nextCursor = posts[len(posts)-1].Id.Hex()
	case dto.OrderByCommentCount:
		var err error
		var previousCompositeKey *string
		if req.PreviousCursor != "" {	
			previousCompositeKey = &req.PreviousCursor
		}
		posts, hasMore, err = service.GetPostService().FetchPostsByCompositCursor(ctx, req.Limit, previousCompositeKey)
		if err != nil {
			hlog.CtxErrorf(ctx, "[HandleGetPost] error fetching posts: %v", err)
			return nil, err
		}
		nextCursor = posts[len(posts)-1].CompositeKey
	}

	var responsePosts []*dto.Post
	for _, post := range posts {
		dtoPost, err := mapDaoPostToDtoPost(ctx, post)
		if err != nil {
			return nil, err
		}
		responsePosts = append(responsePosts, dtoPost)
	}


	var resp = &dto.FetchPostsResp{
		Posts:          responsePosts,
		HasMore:        hasMore,
		PreviousCursor: req.PreviousCursor,
		NextCursor:     nextCursor,
	}

	return resp, nil
}

func mapDaoPostToDtoPost(ctx context.Context, daoPost *dao.Post) (*dto.Post, error) {
	user, err := service.GetUserService().GetUserById(daoPost.Creator)
	if err != nil {
		hlog.CtxErrorf(ctx, "Error fetching user info: %v", err)
		return nil, err
	}
	// Assemble image id and url
	imageID:= filepath.Base(daoPost.ProcessedImagePath)
	imageURL, err := service.GetImageService().GetPublicImageURL(ctx, daoPost.ProcessedImagePath)
	if err != nil {
		hlog.CtxErrorf(ctx, "Error fetching image URL: %v", err)
		return nil, err
	}
	return &dto.Post{
		Id:                     daoPost.Id.Hex(),
		CreatedAtMilli:         daoPost.CreatedAtMilli,
		Content:                daoPost.Content,
		CommentCount:           int(daoPost.CommentCount),
		RecentComments:         mapDaoCommentsToDtoComments(daoPost.RecentComments),
		RecentCommentedAtMilli: daoPost.LastCommentAtMilli,
		CreatorId:              user.Id.Hex(),
		CreatorName:            user.Name,
		ImageId:                imageID,
		ImageURL:               imageURL,
	}, nil
}

func mapDaoCommentsToDtoComments(daoComments []*dao.Comment) []*dto.Comment {
	var dtoComments []*dto.Comment
	for _, daoComment := range daoComments {
		dtoComments = append(dtoComments, &dto.Comment{
			Id:             daoComment.Id.Hex(),
			CreatedAtMilli: daoComment.CreatedAtMilli,
			Content:        daoComment.Content,
			CreatorId:      daoComment.Creator.Hex(),
			CreatorName:    daoComment.CreatorName,
		})
	}
	return dtoComments
}