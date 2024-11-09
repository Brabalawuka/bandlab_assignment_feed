package dto

type GetPresignedURLRequest struct {
	FileName string `query:"fileName" vd:"len($)>0"` // File name
	FileType string `query:"fileType" vd:"len($)>0"` // File type
	FileSize int64  `query:"fileSize" vd:"$>0"` // File size
	UserId   string `header:"userId" vd:"len($)>0"` // User ID
}

// GetPresignedURLResponse represents the response structure for getting a presigned URL
type GetPresignedURLResponse struct {
	ImagePath   string `json:"imagePath"` // Image path e.g "/original/image.jpg"
	URL         string `json:"url"`       // Presigned URL
	ExpiresAtUnix int64  `json:"expiresAtUnix"` // Expiration time in Unix timestamp
}
