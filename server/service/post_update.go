package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/util"
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdatePostStatusAndImageURL updates the image URL of a post and its status
func (s *PostServiceImpl) UpdatePostStatusAndImagePath(ctx context.Context, postId string, imagePath string) error {
	collection := s.mongoClient.Collection(s.mongoCollection)

	// Update the post with the new image path and status
	update := bson.M{
		"$set": bson.M{
			"processedImagePath": imagePath,
			"status":             dao.StatusPosted,
		},
	}

	objId, err := primitive.ObjectIDFromHex(postId)
	if err != nil {
		return fmt.Errorf("invalid post Id: %v, when updating post status and image path", err)
	}

	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, update)
	if result.Err() != nil {
		return fmt.Errorf("mongo error updating post status and image path: %v", result.Err())
	}

	var dao *dao.Post
	if err := result.Decode(&dao); err != nil {
		return fmt.Errorf("failed to decode post after updating post status and image path: %v", err)
	}

	hlog.CtxDebugf(ctx, "Successfully updated post status and image path, postId: %s, imagePath: %s", util.ToJsonString(dao), util.ToJsonString(imagePath))

	return nil
}

// UpdateNewPostComments updates the comments of a post
// If the comment is new, it will be added to the post comment count and recent comments
// If the comment is deleted, it will be removed from the post comment count
func (s *PostServiceImpl) UpdatePostComments(ctx context.Context, postId primitive.ObjectID, comment *dao.Comment, oldPost *dao.Post) (err error) {
	if oldPost == nil {
		oldPost, err = s.FetchPostDaoById(ctx, postId)
		if err != nil {
			hlog.CtxErrorf(ctx, "error fetching post, id: %s, error: %v", postId, err)
			return errs.ErrInternalServer
		}
		if oldPost == nil {
			hlog.CtxErrorf(ctx, "post not found, id: %s", postId)
			return errs.ErrPostNotFound
		}	
	}
	// update the post comments info
	var(
		newCommentCount int32		
		newRecentComments []*dao.Comment
		now = time.Now()
		newLastCommentAtMilli int64
	)

	if comment.Status == dao.CommentStatusDeleted {
		newCommentCount = oldPost.CommentCount - 1
		if newCommentCount < 0 { // this should not happen but just in case
			hlog.CtxWarnf(ctx, "[UpdatePostComments] comment count is negative, postId: %s, commentId: %s", postId, comment.Id)
			newCommentCount = 0
		}
		newRecentComments = filterComments(oldPost.RecentComments, comment.Id) // filter out the deleted comment
		newLastCommentAtMilli = oldPost.LastCommentAtMilli
	} else {
		newCommentCount = oldPost.CommentCount + 1
		newRecentComments = append([]*dao.Comment{comment}, filterComments(oldPost.RecentComments, comment.Id)...)
		if len(newRecentComments) > s.recentCommentsCount {
			newRecentComments = newRecentComments[:s.recentCommentsCount]
		}
		newLastCommentAtMilli = now.UnixMilli()
	}
	// new composite key for the post cursor
	var newCompositeKey string = util.GenerateCompositeKey(int32(newCommentCount), now, oldPost.Id)

	collection := s.mongoClient.Collection("posts")

	// Update the post with the new comment count and recent comments
	update := bson.M{
		"$set": bson.M{
			"commentCount":       newCommentCount,
			"recentComments":     newRecentComments,
			"lastCommentAtMilli": newLastCommentAtMilli,
			"compositeKey":       newCompositeKey,
			"version":            oldPost.Version + 1,
		},
	}

	old := collection.FindOneAndUpdate(ctx, bson.M{"_id": postId, "version": oldPost.Version}, update)
	if old.Err() != nil {
		if old.Err() == mongo.ErrNoDocuments {
			return errs.ErrPostWithVersionNotFound
		}
		return fmt.Errorf("mongo error updating post comments: %v", old.Err())
	}

	return nil
}

// Helper function to filter out a specific comment by ID
func filterComments(comments []*dao.Comment, excludeID primitive.ObjectID) []*dao.Comment {
	var filtered = make([]*dao.Comment, 0, len(comments))
	for _, c := range comments {
		if c.Id != excludeID {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
