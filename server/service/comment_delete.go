package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeleteComment deletes a comment from a post
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.DeleteCommentResp, error) {
	commentId, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		hlog.CtxErrorf(ctx, "[DeleteComment] invalid comment Id, id: %s, error: %v", req.CommentId, err)
		return nil, errs.ErrInvalidRequest
	}

	comment, err := s.fetchComment(ctx, commentId)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		hlog.CtxErrorf(ctx, "[DeleteComment] comment not found, id: %s", commentId.Hex())
		return nil, errs.ErrCommentNotFound
	}
	if !s.isCommentDeletable(comment, req.UserId) {
		return nil, errs.ErrCommentOperationNotAllowed
	}

	softDeletedComment, err := s.softDeleteComment(ctx, commentId)
	if err != nil {
		hlog.CtxErrorf(ctx, "[DeleteComment] failed to mark comment as deleted, id: %s, error: %v", commentId.Hex(), err)
		return nil, errs.ErrInternalServer
	}

	// update the post's comment count and recent comments asynchronously
	// TODO: this is a mock function, in the future, we can use a message queue to update the post because
	// - an async function to update the post num may cause inconsistency
	// - a retry of max 10 times with 100ms interval is used to handle the race condition
	s.updatePostCommentsInfoAsync(ctx, comment.PostId, softDeletedComment, nil)

	return &dto.DeleteCommentResp{Id: commentId.Hex()}, nil
}

func (s *CommentServiceImpl) fetchComment(ctx context.Context, commentId primitive.ObjectID) (*dao.Comment, error) {
	collection := s.mongoClient.Collection(s.mongoCollection)
	var comment dao.Comment
	err := collection.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			hlog.CtxErrorf(ctx, "comment not found, id: %s, error: %v", commentId.Hex(), err)
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (s *CommentServiceImpl) softDeleteComment(ctx context.Context, commentId primitive.ObjectID) (*dao.Comment, error) {
	collection := s.mongoClient.Collection(s.mongoCollection)
	update := bson.M{"$set": bson.M{"status": dao.CommentStatusDeleted}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedComment dao.Comment
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": commentId}, update, opts).Decode(&updatedComment)
	if err != nil {
		hlog.CtxErrorf(ctx, "[softDeleteComment] failed to mark comment as deleted, id: %s, error: %v", commentId.Hex(), err)
		return nil, err
	}
	return &updatedComment, nil
}

func (s *CommentServiceImpl) isCommentDeletable(comment *dao.Comment, userId string) bool {
	if comment.Status != dao.CommentStatusPosted {
		hlog.Errorf("[isCommentDeletable] comment not allowed to delete, id: %s, status: %s", comment.Creator.Hex(), comment.Status)
		return false
	}
	if comment.Creator.Hex() != userId {
		hlog.Errorf("[isCommentDeletable] comment not allowed to delete, id: %s, user id: %s", comment.Creator.Hex(), userId)
		return false
	}
	return true
}
