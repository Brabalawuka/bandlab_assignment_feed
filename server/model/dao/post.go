package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostStatus represents the status of a post
type PostStatus string

const (
	StatusPending PostStatus = "PENDING" // Post is pending since the image resizing is not done
	StatusPosted  PostStatus = "POSTED"  // Post is posted	
	StatusDeleted PostStatus = "DELETED" // TODO: add post soft delete
	// TODO: Add failure status if async image resizing fails
)

// Post represents a post in the database
type Post struct {
	Id                 primitive.ObjectID `bson:"_id,omitempty"` // Post ID
	Content            string             `bson:"content"`       // Post content	
	Creator            primitive.ObjectID `bson:"creator"`       // Post creator ID
	CreatedAtMilli     int64              `bson:"createdAtMilli"` // Post created time in milliseconds	
	Status             PostStatus         `bson:"status"`        // Post status
	CommentCount       int32              `bson:"commentCount"`  // Comment count
	LastCommentAtMilli int64              `bson:"lastCommentAtMilli"` // Last comment time in milliseconds
	OriginalImagePath  string             `bson:"originalImagePath"`  // Original image file path e.g "/original/xxxx.jpg"
	ProcessedImagePath string             `bson:"processedImagePath"` // Processed image file path e.g "/processed/eeeee.jpg"	
	CompositeKey       string             `bson:"compositeKey"`       // Composite key refer to util.CompositeKey
	RecentComments     []*Comment         `bson:"recentComments"`     // Recent comments
	Version            int64              `bson:"version"`            // Version
}


