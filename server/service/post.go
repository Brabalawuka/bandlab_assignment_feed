package service

import (
	"bandlab_feed_server/config"
	"bandlab_feed_server/dal/mongodb"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service/mocks"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
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
//go:generate mockgen -destination=./mocks/mock_post_service.go -package=mocks -source=./post.go
type PostService interface {
	CreatePost(ctx context.Context, req *dto.CreatePostReq) (*dto.CreatePostResp, error)
	HasImage(req *dto.CreatePostReq) bool
	FetchPostDaoById(ctx context.Context, postId primitive.ObjectID) (*dao.Post, error)
	FetchPostsByPostIDCursor(ctx context.Context, limit int64, previousPostId *primitive.ObjectID) (posts []*dao.Post, hasMore bool, err error)
	FetchPostsByCompositCursor(ctx context.Context, limit int64, previousCompositKey *string) (posts []*dao.Post, hasMore bool, err error) 
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
			recentCommentsCount: config.AppConfig.PostRecentCommentsCount,
			mongoClient:     mongoClient,
		}
	})
}

// GetPostService returns the initialized post service
func GetPostService() PostService {
	return postSrv
}

// SetMockPostService For unit testing purpose only
func SetMockPostService(ctrl *gomock.Controller) *mocks.MockPostService {
	mocks := mocks.NewMockPostService(ctrl)
	postSrv = mocks
	return mocks
}

// PostServiceImpl is the implementation of PostService
type PostServiceImpl struct {
	mongoCollection string
	recentCommentsCount int
	mongoClient     mongodb.MongoService
}
