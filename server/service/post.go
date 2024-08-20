package service

import (
	"bandlab_feed_server/dal/mongodb"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents a post in the system
type Post struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	CreatorId      primitive.ObjectID `bson:"creatorId"`
	Content        string             `bson:"content"`
	ImageFileName  string             `bson:"imageFileName"`
	Status         string             `bson:"status"`
	CreatedAtMilli int64              `bson:"createdAt"`
}

// PostService defines the interface for post operations
type PostService interface {
	CreatePost(ctx context.Context, req *dto.CreatePostReq) (*dto.CreatePostResp, error)
	HasImage(req *dto.CreatePostReq) bool
	GetPostDaoById(ctx context.Context, postId primitive.ObjectID) (*dao.Post, error)
	UpdatePostStatusAndImagePath(ctx context.Context, postId string, imagePath string) error
	UpdatePostComments(ctx context.Context, postId primitive.ObjectID, comment *dao.Comment, oldPost *dao.Post) error
}

var (
	postOnce sync.Once
	postSrv  PostService
)

// InitPostService initializes the post service
func InitPostService() {
	postOnce.Do(func() {
		mongoClient := mongodb.GetMongoService()
		if mongoClient == nil {
			panic("mongo client is not initialized")
		}
		postSrv = &PostServiceImpl{
			mongoCollection: "posts",
			mongoClient:     mongoClient,
		}
	})
}

// GetPostService returns the initialized post service
func GetPostService() PostService {
	return postSrv
}

// PostServiceImpl is the implementation of PostService
type PostServiceImpl struct {
	mongoCollection string
	mongoClient     mongodb.MongoService
}
