package response

import "bandlab_feed_server/common/errs"

// GeneralResponse is a generic response structure for all API responses
type GeneralResponse[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data,omitempty"`
}

// NewGeneralResponse creates a new GeneralResponse
func NewGeneralResponse[T any](code int, message string, data T) GeneralResponse[T] {
    return GeneralResponse[T]{
        Code:    code,
        Message: message,
        Data:    data,
    }
}

// NewSuccessResponse creates a new success response with default code 0
func NewSuccessResponse[T any](data T) GeneralResponse[T] {
    return NewGeneralResponse(0, "Operation successful", data)
}

// NewErrorResponse creates a new error response from an APIError
func NewErrorResponse(err *errs.APIError) GeneralResponse[interface{}] {
    return NewGeneralResponse[interface{}](err.Code, err.Message, nil)
}