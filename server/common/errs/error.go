package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents a custom error with a code and message
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"` // Not exposed in JSON
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("error code: %d, message: %s", e.Code, e.Message)
}

// Is implements the Is interface for errors.Is
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// As implements the As interface for errors.As
func (e *APIError) As(target interface{}) bool {
	t, ok := target.(**APIError)
	if !ok {
		return false
	}
	*t = e
	return true
}

// NewBizError creates a new business logic error (HTTP 400)
func NewBizError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message, HTTPStatus: http.StatusBadRequest}
}

// NewInternalError creates a new internal server error (HTTP 500)
func NewInternalError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message, HTTPStatus: http.StatusInternalServerError}
}

// Some predefined error codes (you can expand this list)
const (
	ErrCodeInvalidRequest   = 1001
	ErrCodeUnauthorized     = 1002
	ErrCodeResourceNotFound = 1003
	ErrCodeInvalidInput     = 1004
	ErrCodeUnknownError     = 5000
	ErrCodeInternalError    = 5001

	
	ErrCodeR2ImageNotFound        = 3001
	ErrCodeR2ImageUploadFailed    = 3002
	ErrCodeR2ImageProcessFailed   = 3003
	ErrCodeR2ImageDownLoadFailed  = 3004
	ErrCodeInvalidContentType     = 3005
	ErrCodeImageExistsCheckFailed = 3008
	ErrCodePostNotFound           = 3300
	ErrCodePostWithVersionNotFound = 3301
	ErrCodeCommentNotAllowed      = 3100
	ErrCodeUserNotFound           = 3200
	ErrCodePostFetchFailed        = 3201
	ErrCodePostDataDecodeFailed   = 3202
)

// Predefined errors
var (
	ErrInvalidRequest         = NewBizError(ErrCodeInvalidInput, "Invalid Request")
	ErrUnauthorized           = NewBizError(ErrCodeUnauthorized, "Unauthorized")
	ErrResourceNotFound       = NewBizError(ErrCodeResourceNotFound, "Resource not found")
	ErrInvalidInput           = NewBizError(ErrCodeInvalidInput, "Invalid input")
	ErrInternalServer         = NewInternalError(ErrCodeInternalError, "Internal server error")
	ErrR2ImageDownLoadFailed  = NewInternalError(ErrCodeR2ImageDownLoadFailed, "Failed to download image")
	ErrR2ImageNotFound        = NewBizError(ErrCodeR2ImageNotFound, "Image not found")
	ErrR2ImageUploadFailed    = NewInternalError(ErrCodeR2ImageUploadFailed, "Failed to upload image")
	ErrUserNotFound           = NewBizError(ErrCodeUserNotFound, "User not found")
	ErrInvalidContentType     = NewBizError(ErrCodeInvalidContentType, "Invalid presign content type")
	ErrImageExistsCheckFailed = NewInternalError(ErrCodeImageExistsCheckFailed, "Failed to check if image exists in R2")
	ErrCommentNotAllowed      = NewBizError(ErrCodeCommentNotAllowed, "Comment not allowed")
	ErrPostNotFound           = NewBizError(ErrCodePostNotFound, "Post not found")
	ErrPostWithVersionNotFound = NewBizError(ErrCodePostWithVersionNotFound, "Post with version not found")
	ErrPostFetchFailed        = NewInternalError(ErrCodePostFetchFailed, "Failed to fetch posts")
	ErrPostDataDecodeFailed   = NewInternalError(ErrCodePostDataDecodeFailed, "Failed to decode posts")
)

// GetAPIError converts an error to APIError if possible, or returns a generic internal server error
func GetAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return ErrInternalServer
}
