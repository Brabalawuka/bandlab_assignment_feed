package main

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/common/response"
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// GenericHandlerFunc is a generic type for handler functions
type GenericHandlerFunc[Req any, Resp any] func(context.Context, *Req) (*Resp, error)

// WrapHandler is a generic wrapper for handler functions
// - It binds and validates the request
// - It calls the handler function
// - It handles errors
func WrapHandler[Req any, Resp any](handlerFunc GenericHandlerFunc[Req, Resp]) app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        var req Req
        if err := ctx.BindAndValidate(&req); err != nil {
            hlog.CtxErrorf(c, "[WrapHandler] error binding and validating request: %v", err)
            apiErr := errs.ErrInvalidInput
            ctx.JSON(apiErr.HTTPStatus, response.NewErrorResponse(apiErr))
            return
        }
        var resp *Resp
        resp, err := handlerFunc(c, &req)
        if err != nil {
            var apiErr *errs.APIError
            if errors.As(err, &apiErr) {
                ctx.JSON(apiErr.HTTPStatus, response.NewErrorResponse(apiErr))
            } else {
                // If it's not an APIError, treat it as an unknown error
                internalErr := errs.NewInternalError(errs.ErrCodeUnknownError, err.Error())
                ctx.JSON(internalErr.HTTPStatus, response.NewErrorResponse(internalErr))
            }
            return
        }

        ctx.JSON(200, response.NewSuccessResponse(resp))
    }
}