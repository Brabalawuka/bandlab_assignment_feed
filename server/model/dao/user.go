package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user in the system
type User struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}