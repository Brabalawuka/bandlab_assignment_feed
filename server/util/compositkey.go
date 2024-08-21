package util

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateCompositeKey generates a composite key for a post
// - This is the cursor for the pagination as it is the combination of the comment count, last comment at and the post id
// - It is to an extent unique for MVP version
// - Composite Key Structure:
// 		- 4 bytes: commentCount (uint32, 4.1Billion Comments)
// 		- 4 bytes: lastCommentAt (unix timestamp, till 2109)
// 		- 4 bytes: first 4 bytes of ObjectId (unix timestamp, till 2109)
func GenerateCompositeKey(commentCount int32, lastCommentAt time.Time, postId primitive.ObjectID) string {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint32(buf[0:4], uint32(commentCount))
	binary.BigEndian.PutUint32(buf[4:8], uint32(lastCommentAt.Unix()))
	copy(buf[8:12], postId[:4]) // Assuming Id is a [12]byte
	return base64.StdEncoding.EncodeToString(buf)
}
