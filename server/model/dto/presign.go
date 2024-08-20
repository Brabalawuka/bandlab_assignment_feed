package dto

type GetPresignedURLRequest struct {
	FileName string `query:"fileName" vd:"len($)>0"`
	FileType string `query:"fileType" vd:"len($)>0"`
	FileSize int64  `query:"filesize" vd:"$>0"`
	UserId   string `header:"userId" vd:"len($)>0"`
}

// GetPresignedURLResponse represents the response structure for getting a presigned URL
type GetPresignedURLResponse struct {
	ImagePath   string `json:"imagePath"`
	URL       string `json:"url"`
	ExpiresIn int    `json:"expiresIn"`
}
