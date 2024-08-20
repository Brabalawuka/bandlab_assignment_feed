package service

import (
	"bandlab_feed_server/model/dao"
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *PostServiceImpl) GetPostDaoById(ctx context.Context, postId primitive.ObjectID) (*dao.Post, error) {
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
