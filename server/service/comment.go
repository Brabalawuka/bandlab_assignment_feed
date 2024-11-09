package service

import (
	"bandlab_feed_server/dal/mongodb"
	"bandlab_feed_server/model/dto"
	"context"
	"sync"
)

type CommentService interface {
	CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error)
	DeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.DeleteCommentResp, error)
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
