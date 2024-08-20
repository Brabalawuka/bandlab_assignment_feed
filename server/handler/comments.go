package handler

import (
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service"
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

// CreateComment handles the creation of a new comment on a post
func CreateComment(ctx context.Context, req *dto.CreateCommentReq) (*dto.CreateCommentResp, error) {

	commentService := service.GetCommentService()
	resp, err := commentService.CreateComment(ctx, req)
	if err != nil {
		return nil, err
    }
	return resp, nil
}

// DeleteComment handles the deletion of a comment
func DeleteComment(c context.Context, ctx *app.RequestContext) {
    // 处理删除评论逻辑
    ctx.JSON(http.StatusOK, map[string]string{"message": "Comment deleted"})
}