package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostStatus represents the status of a post
type PostStatus string

const (
	StatusPending PostStatus = "PENDING"
	StatusPosted  PostStatus = "POSTED"
	StatusDeleted PostStatus = "DELETED" // TODO: add post soft delete
	// TODO: Add failure status if async image resizing fails
)

// Post represents a post in the database
type Post struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"`
	Content            string             `bson:"content"`
	Creator            primitive.ObjectID `bson:"creator"`
	CreatedAtMilli     int64              `bson:"createdAtMilli"`
	Status             PostStatus         `bson:"status"`
	CommentCount       int32              `bson:"commentCount"`
	LastCommentAtMilli int64              `bson:"lastCommentAtMilli"`
	OriginalImagePath  string             `bson:"originalImagePath"`
	ProcessedImagePath string             `bson:"processedImagePath"`
	CompositeKey       string             `bson:"compositeKey"`
	RecentComments     []*Comment         `bson:"recentComments"`
	Version            int64              `bson:"version"`
}
