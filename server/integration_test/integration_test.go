package integration_test

import (
	"bandlab_feed_server/common/response"
	"bandlab_feed_server/model/dao"
	"bandlab_feed_server/model/dto"
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/require"
)

var (
	aliceId   = "507f1f77bcf86cd799439011"
	bobId     = "507f1f77bcf86cd799439012"
	charlieId = "507f1f77bcf86cd799439013"
)

var (
	uploadImagePresignUrl string   = "https://bandlab-assignment.bf69da5249b63731ad79545d0095e8db.r2.cloudflarestorage.com/original/66c5e825c0ff015831d181fe.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=7ff1f4515c58056eca10e76c387fbbb9%2F20240821%2Fauto%2Fs3%2Faws4_request&X-Amz-Date=20240821T131413Z&X-Amz-Expires=400&X-Amz-SignedHeaders=content-length%3Bcontent-type%3Bhost&x-id=PutObject&X-Amz-Signature=294acab190d9e5278b48da9cbba2d69f67666c91cd419edda8d8ec7dc930a550"
	uploadImagePath       string   = "original/66c5e825c0ff015831d181fe.jpg"
	uploadURLExpires      int64    = 0
	createdPostID         []string = make([]string, 0)
	createdImageName      string
)

// This is an integration test for the feed service
// This can only be run once as the posts are alr created in the first run
// To run more times you need to comment out the Step 2
func TestIntegration(t *testing.T) {
	// Step 1, Ping Feed Service Operation 
	TestPingFeedService(t)

	// Step 2, test Fetch Post Operation, should return zero post if the mongo is just inited
	//TestFetchPosts(t)

	// Step 3. try to get presigned url of the test_pic_1.jpg in current folder
	TestGetPreSignedURL(t)

	// Step 4, try to create 3 post using the image just uploaded
	TestCreate3Post(t)

	hlog.Infof("Sleeping for 3 seconds to wait for image processing")
	time.Sleep(3 * time.Second) // wait for the image to be processed

	// Step 5, test Fetch Post Operation, should return more posts with last posts same as created post
	// verify the image url is processed same as the image id
	TestFetchPostsInDefaultOrder(t)

	// Step 6, test Comment on post 1, which is the first post created
	TestCommentOnPost1(t)

	time.Sleep(1 * time.Second) // wait for order to be updated

	// Step 7, test Fetch Post Operation, should return posts with more comments ordered in front
	TestFetchPostsInCommentCountOrder(t)
}

func TestPingFeedService(t *testing.T) {
	req := &protocol.Request{}
	req.SetMethod(consts.MethodGet)
	req.SetRequestURI("http://0.0.0.0:8010/ping")

	resp := &protocol.Response{}
	err := hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to ping feed service")
	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

	hlog.Infof("\n-----Step 1: Ping feed service success----\n\n\n")
}

func TestFetchPosts(t *testing.T) {
	req := &protocol.Request{}
	req.SetMethod(consts.MethodGet)

	query := url.Values{}
	query.Set("limit", "10")
	query.Set("orderBy", string(dto.OrderByPostID))
	req.SetRequestURI("http://0.0.0.0:8010/v1/api/posts?" + query.Encode())
	req.SetHeader("userId", aliceId)

	resp := &protocol.Response{}
	err := hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to fetch posts")

	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")
	var fetchPostResp response.GeneralResponse[*dto.FetchPostsResp]
	err = json.Unmarshal(resp.Body(), &fetchPostResp)
	require.NoError(t, err, "Failed to unmarshal response")
	require.Equal(t, 0, fetchPostResp.Code, "Expected code 0")
	require.Equal(t, 0, len(fetchPostResp.Data.Posts), "Expected 0 posts")
	require.Equal(t, false, fetchPostResp.Data.HasMore, "Expected no more posts")

	hlog.Infof("\n-----Step 2: Fetch posts success----\n\n\n")
}

func TestGetPreSignedURL(t *testing.T) {
	// Read the local image file
	filePath, _ := os.Getwd()
	imagePath := filePath + "/test_pic_1.jpg"
	file, err := os.Open(imagePath)
	require.NoError(t, err, "Failed to open image file")
	defer file.Close()

	fileInfo, err := file.Stat()
	require.NoError(t, err, "Failed to get file info")

	fileSize := fileInfo.Size()
	require.NotEqual(t, 0, fileInfo.Size(), "File size should not be 0")
	fileBuffer := make([]byte, fileSize)
	_, err = file.Read(fileBuffer)
	require.NoError(t, err, "Failed to read file")

	mimeType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))

	// Query the presigned url
	req := &protocol.Request{}
	req.SetMethod(consts.MethodGet)

	query := url.Values{}
	query.Set("fileSize", fmt.Sprintf("%d", fileSize))
	query.Set("fileType", mimeType)
	query.Set("fileName", fileInfo.Name())
	req.SetRequestURI("http://0.0.0.0:8010/v1/api/posts/image-presign?" + query.Encode())
	req.SetHeader("userId", aliceId)

	resp := &protocol.Response{}
	err = hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to get presigned url")
	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

	var presignResp response.GeneralResponse[*dto.GetPresignedURLResponse]
	err = json.Unmarshal(resp.Body(), &presignResp)
	require.NoError(t, err, "Failed to unmarshal response")
	require.Equal(t, 0, presignResp.Code, "Expected code 0")
	require.NotEmpty(t, presignResp.Data.URL, "Expected non-empty URL")

	// Receive the presigned url and image id
	uploadImagePath = presignResp.Data.ImagePath
	uploadImagePresignUrl = presignResp.Data.URL
	uploadURLExpires = presignResp.Data.ExpiresAtUnix
	createdImageName = filepath.Base(uploadImagePath)

	// upload the image to the presigned url
	uploadReq := &protocol.Request{}
	uploadReq.SetMethod(consts.MethodPut)
	uploadReq.SetRequestURI(uploadImagePresignUrl)
	uploadReq.SetBody(fileBuffer)
	uploadReq.SetHeader("Content-Type", mimeType)
	uploadReq.SetHeader("Content-Length", fmt.Sprintf("%d", fileSize))
	hlog.Infof("upload image presigned url: %s", uploadImagePresignUrl)
	uploadResp := &protocol.Response{}
	err = hertzClient.Do(context.Background(), uploadReq, uploadResp)
	require.NoError(t, err, "Failed to upload image")
	require.Equal(t, consts.StatusOK, uploadResp.StatusCode(), "Expected status code 200")

	hlog.Infof("\n-----Step 3: Get presigned url success, and upload image success----\n")
	hlog.Infof("presigned url: %s, image path: %s, expire time: %d", uploadImagePresignUrl, uploadImagePath, uploadURLExpires)
	hlog.Infof("created image name: %s", createdImageName)
	hlog.Infof("\n-------------------------------------------------\n\n\n")

}

func TestCreate3Post(t *testing.T) {
	// Create post request payload
	for i := 0; i < 3; i++ {
		postPayload := dto.CreatePostReq{
			ImageFilePath: uploadImagePath,
			Content:       fmt.Sprintf("Test post content %d, %s", i+1, t.Name()),
		}

		// Marshal the payload to JSON
		postPayloadBytes, err := json.Marshal(postPayload)
		require.NoError(t, err, "Failed to marshal post payload")

		// Create the request
		req := &protocol.Request{}
		req.SetMethod(consts.MethodPost)
		req.SetRequestURI("http://0.0.0.0:8010/v1/api/posts")
		req.SetBody(postPayloadBytes)
		req.SetHeader("Content-Type", "application/json")
		req.SetHeader("userId", aliceId)

		// Send the request
		resp := &protocol.Response{}
		err = hertzClient.Do(context.Background(), req, resp)
		require.NoError(t, err, "Failed to create post")
		require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

		// Unmarshal the response
		var createPostResp response.GeneralResponse[*dto.CreatePostResp]
		err = json.Unmarshal(resp.Body(), &createPostResp)
		require.NoError(t, err, "Failed to unmarshal response")
		require.Equal(t, 0, createPostResp.Code, "Expected code 0")
		require.Equal(t, createPostResp.Data.Content, postPayload.Content, "Expected content to match")
		require.NotNil(t, createPostResp.Data.CreatorId, aliceId)
		require.Equal(t, createPostResp.Data.Status, string(dao.StatusPending), "Expected status to be pending")
		createdPostID = append(createdPostID, createPostResp.Data.Id)

		hlog.Infof("\n-----Step 4: Create post No%d success----\n", i+1)
		hlog.Infof("post id: %s, content: %s", createPostResp.Data.Id, createPostResp.Data.Content)
		hlog.Infof("\n-------------------------------------------------\n\n\n")
	}

}

func TestFetchPostsInDefaultOrder(t *testing.T) {
	req := &protocol.Request{}
	req.SetMethod(consts.MethodGet)

	query := url.Values{}
	query.Set("limit", "50")
	query.Set("orderBy", string(dto.OrderByPostID))
	req.SetRequestURI("http://0.0.0.0:8010/v1/api/posts?" + query.Encode())
	req.SetHeader("userId", aliceId)

	resp := &protocol.Response{}
	err := hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to fetch posts")
	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

	var fetchPostResp response.GeneralResponse[*dto.FetchPostsResp]
	err = json.Unmarshal(resp.Body(), &fetchPostResp)
	require.NoError(t, err, "Failed to unmarshal response")
	lastPost := fetchPostResp.Data.Posts[0] // last post is the first fetch record
	require.Equal(t, 0, fetchPostResp.Code, "Expected code 0")
	require.Greater(t, len(fetchPostResp.Data.Posts), 0, "Greater than 0 posts")
	require.Equal(t, false, fetchPostResp.Data.HasMore, "Expected no more posts")
	require.Equal(t, createdPostID[len(createdPostID)-1], lastPost.Id, "Expected post id to match")
	require.Equal(t, "https://pub-e7ae7f4305084b3ea6e32696f803a332.r2.dev/processed/"+createdImageName, lastPost.ImageURL, "Expected content to match")

	hlog.Infof("\n-----Step 5: Fetch posts success----\n")
	hlog.Infof("last post id: %s, image url: %s", lastPost.Id, lastPost.ImageURL)
	hlog.Infof("\n-------------------------------------------------\n\n\n")
}

func TestCommentOnPost1(t *testing.T) {
	// Create comment request payload
	commentPayload := dto.CreateCommentReq{
		Content: "Test comment1 content1 lloremipsum",
		PostId:  createdPostID[0],
	}

	// Marshal the payload to JSON
	commentPayloadBytes, err := json.Marshal(commentPayload)
	require.NoError(t, err, "Failed to marshal comment payload")

	// Create the request
	req := &protocol.Request{}
	req.SetMethod(consts.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://0.0.0.0:8010/v1/api/posts/%s/comments", createdPostID[0]))
	req.SetBody(commentPayloadBytes)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("userId", bobId)
	var resp = &protocol.Response{}
	err = hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to create comment")
	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

	// Unmarshal the response
	var createCommentResp response.GeneralResponse[*dto.CreateCommentResp]
	err = json.Unmarshal(resp.Body(), &createCommentResp)
	require.NoError(t, err, "Failed to unmarshal response")
	require.Equal(t, 0, createCommentResp.Code, "Expected code 0")
	require.Equal(t, createdPostID[0], createCommentResp.Data.PostId, "Expected post id to match")

	hlog.Infof("\n-----Step 6: Comment on post 1 success----\n\n")
	hlog.Infof("comment id: %s, post id: %s", createCommentResp.Data.Id, createCommentResp.Data.PostId)
	hlog.Infof("\n-------------------------------------------------\n\n\n")
}

func TestFetchPostsInCommentCountOrder(t *testing.T) {
	req := &protocol.Request{}
	req.SetMethod(consts.MethodGet)

	query := url.Values{}
	query.Set("limit", "50")
	query.Set("orderBy", string(dto.OrderByCommentCount))
	req.SetRequestURI("http://0.0.0.0:8010/v1/api/posts?" + query.Encode())
	req.SetHeader("userId", aliceId)

	resp := &protocol.Response{}
	err := hertzClient.Do(context.Background(), req, resp)
	require.NoError(t, err, "Failed to fetch posts")

	require.Equal(t, consts.StatusOK, resp.StatusCode(), "Expected status code 200")

	var fetchPostResp response.GeneralResponse[*dto.FetchPostsResp]
	err = json.Unmarshal(resp.Body(), &fetchPostResp)
	commenttedPost := fetchPostResp.Data.Posts[0] // commented post is the first fetch record
	require.NoError(t, err, "Failed to unmarshal response")
	require.Equal(t, 0, fetchPostResp.Code, "Expected code 0")
	require.Greater(t, len(fetchPostResp.Data.Posts), 0, "Greater than 0 posts")
	require.Equal(t, false, fetchPostResp.Data.HasMore, "Expected no more posts")
	require.Equal(t, createdPostID[0], commenttedPost.Id, "Expected post id to match")
	require.Equal(t, "https://pub-e7ae7f4305084b3ea6e32696f803a332.r2.dev/processed/"+createdImageName, commenttedPost.ImageURL, "Expected content to match")
	hlog.Infof("\n-----Step 7: Fetch posts in comment count order success----\n\n")
	hlog.Infof("commentted post id: %s, image url: %s", commenttedPost.Id, commenttedPost.ImageURL)
	hlog.Infof("\n-------------------------------------------------\n\n\n")
}


