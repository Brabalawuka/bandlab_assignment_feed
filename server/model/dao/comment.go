package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment in the database
type Comment struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`      // Comment ID
	PostId         primitive.ObjectID `bson:"postId"`             // Post ID
	ParentId       primitive.ObjectID `bson:"parentId,omitempty"` // Parent comment ID (optional)
	Content        string             `bson:"content"`            // Comment content
	Status         CommentStatus      `bson:"status"`             // Comment status
	Creator        primitive.ObjectID `bson:"creator"`            // Comment creator ID
	CreatorName    string             `bson:"creatorName"`        // Comment creator name
	CreatedAtMilli int64              `bson:"createdAtMilli"`     // Comment created time in milliseconds
}

// CommentStatus represents the status of a comment
type CommentStatus string

const (
	CommentStatusPosted  CommentStatus = "POSTED"  // Comment is posted
	CommentStatusDeleted CommentStatus = "DELETED" // Comment is deleted
)
