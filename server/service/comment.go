package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/dal/mongodb"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/util/async"
	"context"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommentService interface {
	CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error)
}

// CommentServiceImpl implements the PostCommentService interface
type CommentServiceImpl struct {
	mongoClient     mongodb.MongoService
	mongoCollection string
}

var (
	commentOnce sync.Once
	commentSrv  CommentService
)

// InitCommentService initializes the comment service
func InitCommentService() {
	commentOnce.Do(func() {
		mongoClient := mongodb.GetMongoService()
		if mongoClient == nil {
			panic("mongo client is not initialized")
		}
		commentSrv = &CommentServiceImpl{
			mongoCollection: "comments",
			mongoClient:     mongoClient,
		}
	})
}

// GetCommentService returns the initialized comment service
func GetCommentService() CommentService {
	return commentSrv
}

// PostComment creates a new comment on a post
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error) {
	// Convert UserId and PostId to ObjectId
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		hlog.CtxErrorf(ctx, "invalid user Id, id: %s, error: %v", req.UserId, err)
		return nil, errs.ErrInvalidRequest
	}
	postId, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		hlog.CtxErrorf(ctx, "invalid post Id, id: %s, error: %v", req.PostId, err)
		return nil, errs.ErrInvalidRequest
	}
	// Fetch user info
	user, err := GetUserService().GetUserById(userId)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get user, error: %v", err)
		return nil, errs.ErrInternalServer
	}
	//Fetch the post to verify the post exists
	post, err := GetPostService().FetchPostDaoById(ctx, postId)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get post, error: %v", err)
		return nil, errs.ErrInternalServer
	}
	if post == nil {
		hlog.CtxWarnf(ctx, "post not found, id: %s", postId)
		return nil, errs.ErrPostNotFound
	}
	// if post is not published, return error
	if post.Status != dao.StatusPosted {
		hlog.CtxErrorf(ctx, "post is not published, id: %s", postId)
		return nil, errs.ErrCommentNotAllowed
	}

	// Create a new comment
	createdAt := time.Now().UnixMilli()
	comment := &dao.Comment{
		PostId:         postId,
		Content:        req.Content,
		Status:         "ACTIVE",
		Creator:        user.Id,
		CreatorName:    user.Name, // TODO: Suport dynamic name in the future
		CreatedAtMilli: createdAt,
	}

	// Insert the comment into the database
	comment, err = s.insertComment(ctx, comment)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to insert comment, error: %v", err)
		return nil, errs.ErrInternalServer
	}

	// Asynchronously update the post's comment count and recent comments
	// This is a mock function, in the future, we will use a message queue to update the post
	// - An async function to update the post num may cause inconsistency
	// - A retry of max 10 times with 100ms interval is used to handle the race condition
	async.Go(ctx, "UpdatePostComments", func(ctx context.Context) {
		var err error
		for i := 0; i < 10; i++ { // retry 10 times
			err = GetPostService().UpdatePostComments(ctx, postId, comment, post)
			if err == nil {
				break
			}
			if err != nil && err != errs.ErrPostWithVersionNotFound {
				break
			}
			// optimistic lock, retry after 100ms
			time.Sleep(time.Millisecond * 100)
		}
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update post comments, error: %v", err)
		}
	})

	// Create response
	resp := &dto.CreateCommentResp{
		Id:             comment.Id.Hex(),
		PostId:         comment.PostId.Hex(),
		CreatedAtMilli: createdAt,
	}

	return resp, nil
}

// insertComment inserts a comment into the database
func (s *CommentServiceImpl) insertComment(ctx context.Context, comment *dao.Comment) (*dao.Comment, error) {
	// Insert the comment into the database
	collection := s.mongoClient.Collection(s.mongoCollection)
	result, err := collection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	insertedId, _ := result.InsertedID.(primitive.ObjectID)
	comment.Id = insertedId
	return comment, nil
}
