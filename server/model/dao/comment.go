package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment in the database
type Comment struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	PostId         primitive.ObjectID `bson:"postId"`
	Content        string             `bson:"content"`
	Status         CommentStatus      `bson:"status"`
	Creator        primitive.ObjectID `bson:"creator"`
	CreatedAtMilli int64              `bson:"createdAtMilli"`
}

// CommentStatus represents the status of a comment
type CommentStatus string

const (
	CommentStatusPosted  CommentStatus = "POSTED"
	CommentStatusDeleted CommentStatus = "DELETED"
)
