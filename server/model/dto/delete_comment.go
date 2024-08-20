package dto

// DeleteCommentReq is a struct that represents the DeleteCommentReq DTO.
type DeleteCommentReq struct {
	CommentId string `json:"commentId" vd:"len($)>0"`
	UserId    string `header:"userId" vd:"len($)>0"`
}

// DeleteCommentResp is a struct that represents the DeleteCommentResp DTO.
type DeleteCommentResp struct {
	Id             string `json:"id"`
	DeletedAtMilli int64  `json:"deletedAt"`
}
