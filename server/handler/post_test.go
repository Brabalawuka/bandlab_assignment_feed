package handler

import (
	"bandlab_feed_server/common/errs"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"bandlab_feed_server/service"
	"bandlab_feed_server/service/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestHandleGetPost(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockPostSrv := service.SetMockPostService(ctrl)
	mockUserSrc := service.SetMockUserService(ctrl)
	mockImageSrc := service.SetMockImageService(ctrl)

	tests := []struct {
		name  string
		ctx context.Context
		provide *dto.FetchPostsReq
		wantResp bool
		wantPost  dto.Post
		wantHasMore bool
		wantCursor string
		wantErr error
		stub  func(mockPostSrv *mocks.MockPostService, mockUserSrc *mocks.MockUserService, mockImageSrc *mocks.MockImageService)
	}{
		{
			name: "Invalid Cursor",
			ctx: context.Background(),
			provide: &dto.FetchPostsReq{
				Limit: 10,
				PreviousCursor: "invalid",
				UserId: "507f1f77bcf86cd799439011",
				OrderBy: dto.OrderByPostID,
			},
			wantResp: false,
			wantErr: errs.ErrInvalidRequest,
			stub: func(mockPostSrv *mocks.MockPostService, mockUserSrc *mocks.MockUserService, mockImageSrc *mocks.MockImageService) {
				return
			},
		},
		{
			name: "Fetch Posts By Post ID Error",
			ctx: context.Background(),
			provide: &dto.FetchPostsReq{
				Limit: 10,
				PreviousCursor: "",
				UserId: "507f1f77bcf86cd799439011",
				OrderBy: dto.OrderByPostID,
			},
			wantResp: false,
			wantErr: errs.ErrInternalServer,
			stub: func(mockPostSrv *mocks.MockPostService, mockUserSrc *mocks.MockUserService, mockImageSrc *mocks.MockImageService) {
				mockPostSrv.EXPECT().FetchPostsByPostIDCursor(gomock.Any(), int64(10), nil).Return(nil, false, errs.ErrInternalServer)
			},
		},
		{
			name: "Fetch Posts By Comment Count Success",
			ctx: context.Background(),
			provide: &dto.FetchPostsReq{
				Limit: 10,
				PreviousCursor: "",
				UserId: "507f1f77bcf86cd799439011",
				OrderBy: dto.OrderByCommentCount,
			},
			wantResp: true,
			wantPost: dto.Post{
				Id: "507f1f77bcf86cd799439011",
				CreatedAtMilli: 1000,
				Content: "content",
				CommentCount: 0,
				RecentComments: nil,
				RecentCommentedAtMilli: 0,
				CreatorId: "507f1f77bcf86cd799439011",
				CreatorName: "Alice",
				ImageURL: "mock_url",
			},
			wantHasMore: false,
			wantCursor: "",
			wantErr: nil,
			stub: func(mockPostSrv *mocks.MockPostService, mockUserSrc *mocks.MockUserService, mockImageSrc *mocks.MockImageService) {
				id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				creatorID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				mockPostSrv.EXPECT().FetchPostsByCompositCursor(gomock.Any(), int64(10), nil).Return([]*dao.Post{
					{
						Id: id ,
						CreatedAtMilli: 1000,
						Content: "content",
						CommentCount: 0,
						RecentComments: nil,
						LastCommentAtMilli: 0,
						Creator: creatorID,
						ProcessedImagePath: "",
						CompositeKey: "",
					},
				}, false, nil)
				mockUserSrc.EXPECT().GetUserById(creatorID).Return(&dao.User{
					Id: creatorID,
					Name: "Alice",
				}, nil)
				mockImageSrc.EXPECT().GetPublicImageURL(gomock.Any(), "").Return("mock_url", nil)
			},
		},


		
	
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.stub(mockPostSrv, mockUserSrc, mockImageSrc)
			gotResp, gotErr := HandleGetPost(tt.ctx, tt.provide)
			if tt.wantResp{
				gotPost := *gotResp.Posts[0]
				assert.Equal(t, gotPost, tt.wantPost)
				assert.Equal(t, gotResp.HasMore, tt.wantHasMore)
				assert.Equal(t, gotResp.PreviousCursor, tt.wantCursor)
			} else {
				assert.Nil(t, gotResp)
			}
			assert.Equal(t, gotErr, tt.wantErr)
		})
	}
}
