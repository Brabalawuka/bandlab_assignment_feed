package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/util/async"
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateComment creates a new comment on a post
func (s *CommentServiceImpl) CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error) {
	userId, postId, err := s.validateAndConvertIds(ctx, req)
	if err != nil {
		return nil, err
	}
	// fetch the user info and post info for validation and information
	user, post, err := s.fetchUserAndPost(ctx, userId, postId)
	if err != nil {
		return nil, err
	}

	// create and insert the comment into the database
	comment, err := s.createAndInsertComment(ctx, req, user, post)
	if err != nil {
		return nil, err
	}

	// update the post's comment count and recent comments asynchronously
	// TODO: this is a mock function, in the future, we can use a message queue to update the post because
	// - an async function to update the post num may cause inconsistency
	// - a retry of max 10 times with 100ms interval is used to handle the race condition
	s.updatePostCommentsInfoAsync(ctx, postId, comment, post)

	resp := &dto.CreateCommentResp{
		Id:             comment.Id.Hex(),
		PostId:         comment.PostId.Hex(),
		CreatedAtMilli: comment.CreatedAtMilli,
	}

	return resp, nil
}

func (s *CommentServiceImpl) validateAndConvertIds(ctx context.Context, req *dto.CreateCommentReq) (primitive.ObjectID, primitive.ObjectID, error) {
	userId, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		hlog.CtxErrorf(ctx, "invalid user Id, id: %s, error: %v", req.UserId, err)
		return primitive.NilObjectID, primitive.NilObjectID, errs.ErrInvalidRequest
	}
	postId, err := primitive.ObjectIDFromHex(req.PostId)
	if err != nil {
		hlog.CtxErrorf(ctx, "invalid post Id, id: %s, error: %v", req.PostId, err)
		return primitive.NilObjectID, primitive.NilObjectID, errs.ErrInvalidRequest
	}
	return userId, postId, nil
}

func (s *CommentServiceImpl) fetchUserAndPost(ctx context.Context, userId, postId primitive.ObjectID) (*dao.User, *dao.Post, error) {
	user, err := GetUserService().GetUserById(userId)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get user, error: %v", err)
		return nil, nil, errs.ErrInternalServer
	}

	post, err := GetPostService().FetchPostDaoById(ctx, postId)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get post, error: %v", err)
		return nil, nil, errs.ErrInternalServer
	}
	if post == nil {
		hlog.CtxWarnf(ctx, "post not found, id: %s", postId)
		return nil, nil, errs.ErrPostNotFound
	}
	if post.Status != dao.StatusPosted {
		hlog.CtxErrorf(ctx, "post is not published, id: %s", postId)
		return nil, nil, errs.ErrCommentNotAllowed
	}

	return user, post, nil
}

func (s *CommentServiceImpl) createAndInsertComment(ctx context.Context, req *dto.CreateCommentReq, user *dao.User, post *dao.Post) (*dao.Comment, error) {
	createdAt := time.Now().UnixMilli()
	comment := &dao.Comment{
		PostId:         post.Id,
		Content:        req.Content,
		Status:         dao.CommentStatusPosted,
		Creator:        user.Id,
		CreatorName:    user.Name,
		CreatedAtMilli: createdAt,
	}

	// Insert the comment into the database
	collection := s.mongoClient.Collection(s.mongoCollection)
	result, err := collection.InsertOne(ctx, comment)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to insert comment, error: %v", err)
		return nil, errs.ErrInternalServer
	}
	insertedId, _ := result.InsertedID.(primitive.ObjectID)
	comment.Id = insertedId

	return comment, nil
}

// update the post's comment count and recent comments asynchronously (handles both create and delete)
// TODO: this is a mock function, in the future, we can use a message queue to update the post because
// - an async function to update the post num may cause inconsistency
// - a retry of max 10 times with 100ms interval is used to handle the race condition
func (s *CommentServiceImpl) updatePostCommentsInfoAsync(ctx context.Context, postId primitive.ObjectID, comment *dao.Comment, post *dao.Post) {
	async.Go(ctx, "UpdatePostComments", func(ctx context.Context) {
		var err error
		for i := 0; i < 10; i++ {
			err = GetPostService().UpdatePostComments(ctx, postId, comment, post)
			if err == nil {
				break
			}
			if err != nil && err != errs.ErrPostWithVersionNotFound {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to update post comments, error: %v", err)
		}
	})
}
