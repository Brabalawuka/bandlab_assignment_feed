package handler

import (
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service"
	"context"
)

// Entry point for creating a comment
func CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error) {

	commentService := service.GetCommentService()
	resp, err := commentService.CreateComment(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Entry point for deleting a comment
func DeleteComment(ctx context.Context, req *dto.DeleteCommentReq) (*dto.DeleteCommentResp, error) {
	commentService := service.GetCommentService()
	resp, err := commentService.DeleteComment(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
