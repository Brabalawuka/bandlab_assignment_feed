package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/util"
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreatePost creates a new post
func (s *PostServiceImpl) CreatePost(ctx context.Context, req *dto.CreatePostReq) (*dto.CreatePostResp, error) {
	// Convert UserId to ObjectId
	// TODO: extract to a validator
	creatorId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		hlog.CtxErrorf(ctx, "invalid user Id, id: ", req.UserId, "error: ", err)
		return nil, errs.ErrInvalidRequest
	}
	// validate UserId
	if _, err := userSrv.GetUserById(creatorId); err != nil {
		hlog.CtxErrorf(ctx, "user not found, id: ", req.UserId)
		return nil, errs.ErrUserNotFound
	}
	// validate image if exists
	if req.ImageFilePath != "" {
		// check if image exists in R2
		if exists, err := imageSrv.RawImageExists(ctx, req.ImageFilePath); err != nil {
			hlog.CtxErrorf(ctx, "failed to check if image exists in R2, error: ", err)
			return nil, err
		} else if !exists {
			hlog.CtxErrorf(ctx, "image does not exist in R2, path: ", req.ImageFilePath)
			return nil, errs.ErrR2ImageNotFound
		}
	}
	// validate content length
	if len(req.Content) > 1000 {
		hlog.CtxErrorf(ctx, "content length is too long, length: ", len(req.Content))
		return nil, errs.ErrInvalidRequest
	}

	// Create a new post
	createdAt := time.Now().UnixMilli()
	id := primitive.NewObjectID()
	post := &dao.Post{
		Id:                 id,
		Creator:            creatorId,
		Content:            req.Content,
		CreatedAtMilli:     createdAt,
		CommentCount:       0,
		LastCommentAtMilli: 0,
		CompositeKey:       util.GenerateCompositeKey(0, time.UnixMilli(createdAt), id), // This is the pagination cursor
		RecentComments:     []*dao.Comment{},
	}
	// set status
	if s.HasImage(req) {
		post.OriginalImagePath = req.ImageFilePath
		post.Status = dao.StatusPending
	} else {
		post.Status = dao.StatusPosted
	}

	// Insert the post into the database
	insertedPost, err := s.insertPost(ctx, post)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to insert post, error: ", err)
		return nil, errs.ErrInternalServer
	}

	// Create response
	resp := &dto.CreatePostResp{
		Id:             insertedPost.Id.Hex(),
		CreatorId:      creatorId.Hex(),
		Content:        insertedPost.Content,
		Status:         string(insertedPost.Status),
		CreatedAtMilli: createdAt,
	}

	return resp, nil
}

func (s *PostServiceImpl) HasImage(req *dto.CreatePostReq) bool {
	return req.ImageFilePath != ""
}

// insertPost inserts a post into the database
func (s *PostServiceImpl) insertPost(ctx context.Context, post *dao.Post) (*dao.Post, error) {
	collection := s.mongoClient.Collection(s.mongoCollection)
	result, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	insertedId, _ := result.InsertedID.(primitive.ObjectID)
	post.Id = insertedId
	return post, nil
}
