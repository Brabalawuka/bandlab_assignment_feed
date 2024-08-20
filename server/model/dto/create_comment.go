package dto

// CreateCommentReq is a struct that represents the CreateCommentReq DTO.
type CreateCommentReq struct {
	Content  string `json:"content" vd:"len($)>0&& len($)<=1000"` //TODO: Dynamic validation of content length
	PostId   string `path:"postId" vd:"len($)>0"`
	ParentId string `json:"parentId"` // TODO: Comments on a comment
	UserId   string `header:"userId" vd:"len($)>0"`
}

// CreateCommentResp is a struct that represents the CreateCommentResp DTO.
type CreateCommentResp struct {
	Id             string `json:"id"`
	PostId         string `json:"postId"`
	CreatedAtMilli int64  `json:"createdAt"`
}
