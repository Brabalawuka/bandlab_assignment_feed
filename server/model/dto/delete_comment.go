package dto

// DeleteCommentReq is a struct that represents the DeleteCommentReq DTO.
type DeleteCommentReq struct {
	CommentId string `json:"commentId" vd:"len($)>0"` // Comment ID
	UserId    string `header:"userId" vd:"len($)>0"` // User ID
}

// DeleteCommentResp is a struct that represents the DeleteCommentResp DTO.
type DeleteCommentResp struct {
	Id             string `json:"id"` // Comment ID
	DeletedAtMilli int64  `json:"deletedAt"` // Comment deleted time in milliseconds
}
