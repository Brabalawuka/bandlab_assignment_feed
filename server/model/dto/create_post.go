package dto

// CreatePostReq is a struct that represents the CreatePostReq DTO.
type CreatePostReq struct {
	Content       string `json:"content" vd:"len($)>0&& len($)<=1000"`
	ImageFilePath string `json:"imageFilePath"`
	UserId        string `header:"userId" vd:"len($)>0"`
}

type CreatePostResp struct {
	Id             string `json:"id"`
	CreatorId      string `json:"creatorId"`
	Content        string `json:"content"`
	Status         string `json:"status"`
	CreatedAtMilli int64  `json:"createdAt"`
}
