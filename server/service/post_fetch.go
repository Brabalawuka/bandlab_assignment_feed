package service

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FetchPostDaoById fetch post by id, if not found, return nil, nil, otherwise return post and nil
func (s *PostServiceImpl) FetchPostDaoById(ctx context.Context, postId primitive.ObjectID) (*dao.Post, error) {
	collection := s.mongoClient.Collection("posts")
	post := &dao.Post{}
	err := collection.FindOne(ctx, bson.M{"_id": postId}).Decode(post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			hlog.CtxWarnf(ctx, "post not found, id: %s", postId)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch post, id: %s, error: %v", postId, err)
	}

	return post, nil
}

func (s *PostServiceImpl) fetchPosts(ctx context.Context, limit int64, filter bson.M, sortKey string) (posts []*dao.Post, hasMore bool, err error) {
	collection := s.mongoClient.Collection("posts")
	cursor, err := collection.Find(
		ctx,
		filter,
		options.Find().SetLimit(limit + 1).SetSort(bson.D{{Key: sortKey, Value: -1}}),
	)
	if err != nil {
		hlog.CtxErrorf(ctx, "[fetchPosts] failed to fetch posts, error: %v", err)
		return nil, false, errs.ErrPostFetchFailed
	}

	if err := cursor.All(ctx, &posts); err != nil {
		hlog.CtxErrorf(ctx, "[fetchPosts] failed to decode posts, error: %v", err)
		return nil, false, errs.ErrPostDataDecodeFailed
	}

	hasMore = len(posts) > int(limit)
	hlog.Debugf("[fetchPosts] hasMore: %v, limit: %v, len(posts): %v", hasMore, limit, len(posts))
	if hasMore {
		posts = posts[:len(posts)-1]
	}

	return posts, hasMore, nil
}

func (s *PostServiceImpl) FetchPostsByPostIDCursor(ctx context.Context, limit int64, previousPostId *primitive.ObjectID) (posts []*dao.Post, hasMore bool, err error) {
	filter := bson.M{
		"status": bson.M{"$eq": dao.StatusPosted},
	}
	if previousPostId != nil {
		filter["_id"] = bson.M{"$lt": *previousPostId}
	}
	return s.fetchPosts(ctx, limit, filter, "_id")
}

func (s *PostServiceImpl) FetchPostsByCompositCursor(ctx context.Context, limit int64, previousCompositKey *string) (posts []*dao.Post, hasMore bool, err error) {
	filter := bson.M{
		"status": bson.M{"$eq": dao.StatusPosted},
	}
	if previousCompositKey != nil {
		filter["compositeKey"] = bson.M{"$lt": *previousCompositKey}
	}
	return s.fetchPosts(ctx, limit, filter, "compositeKey")
}