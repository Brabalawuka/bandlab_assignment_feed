package dto

// CreatePostReq is a struct that represents the CreatePostReq DTO.
type CreatePostReq struct {
	Content       string `json:"content" vd:"len($)>0&& len($)<=1000"` // Post content
	ImageFilePath string `json:"imageFilePath"` // Image file path
	UserId        string `header:"userId" vd:"len($)>0"` // User ID
}

type CreatePostResp struct {
	Id             string `json:"id"` // Post ID
	CreatorId      string `json:"creatorId"` // Post creator ID
	Content        string `json:"content"` // Post content
	Status         string `json:"status"` // Post status
	CreatedAtMilli int64  `json:"createdAt"` // Post created time in milliseconds
}
