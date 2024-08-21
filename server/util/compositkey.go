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
// 		- 8 bytes: last 8 bytes of PostID (ObjectId)
func GenerateCompositeKey(commentCount int32, lastCommentAt time.Time, postID primitive.ObjectID) string {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[0:4], uint32(commentCount))
	binary.BigEndian.PutUint32(buf[4:8], uint32(lastCommentAt.Unix()))
	copy(buf[8:16], postID[4:]) // Assuming Id is a [12]byte
	return base64.StdEncoding.EncodeToString(buf)
}
