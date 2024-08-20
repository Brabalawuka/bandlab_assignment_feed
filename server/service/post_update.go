package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/util"
	"context"
	"fmt"
	"time"

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

	return nil
}

// UpdateNewPostComments updates the comments of a post
// If the comment is new, it will be added to the post comment count and recent comments
// If the comment is deleted, it will be removed from the post comment count
func (s *PostServiceImpl) UpdatePostComments(ctx context.Context, postId primitive.ObjectID, comment *dao.Comment, oldPost *dao.Post) (err error) {
	if oldPost == nil {
		oldPost, err = s.GetPostDaoById(ctx, postId)
		if err != nil {
			return errs.ErrPostNotFound
		}
	}
	var newCommentCount int32
	var newRecentComments []*dao.Comment
	var now time.Time
	var newLastCommentAtMilli int64
	if comment.Status == dao.CommentStatusDeleted {
		newCommentCount = oldPost.CommentCount - 1
		for _, c := range oldPost.RecentComments {
			if c.Id != comment.Id {
				newRecentComments = append(newRecentComments, c)
			}
		}
		newLastCommentAtMilli = oldPost.LastCommentAtMilli
	} else {
		newCommentCount = oldPost.CommentCount + 1
		newRecentComments = append(oldPost.RecentComments, comment)
		newLastCommentAtMilli = now.UnixMilli()
	}
	var newCompositeKey string = util.GenerateCompositeKey(int32(newCommentCount), now, oldPost.Id)

	collection := s.mongoClient.Collection("posts")

	// Update the post with the new comment count and recent comments
	update := bson.M{
		"$set": bson.M{
			"commentCount":       newCommentCount,
			"recentComments":     newRecentComments,
			"lastCommentAtMilli": newLastCommentAtMilli,
			"compositeKey":       newCompositeKey,
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
