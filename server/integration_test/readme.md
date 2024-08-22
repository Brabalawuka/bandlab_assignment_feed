# Feed Service Testing Scripts

This README contains curl commands and a bash script for testing various endpoints of the feed service.

## Prerequisites

- Ensure that the feed service is running on `http://0.0.0.0:8010`
- Install `curl` and `jq` on your system if they're not already available

## Test Scripts

### Step 1: Ping Feed Service

```bash
curl -X GET "http://0.0.0.0:8010/ping"
```

### Step 2: Fetch Posts (Initial State)

```bash
curl -X GET "http://0.0.0.0:8010/v1/api/posts?limit=10&orderBy=post_id" \
     -H "userId: 507f1f77bcf86cd799439011"
```

### Step 3: Get PreSigned URL and Upload Image

```bash
# Get the presigned URL
curl -X GET "http://0.0.0.0:8010/v1/api/posts/image-presign?fileSize={{fileSize}}&fileType={{mimeType}}&fileName={{fileName}}" \
     -H "userId: 507f1f77bcf86cd799439011"

# Upload the image to the presigned URL
curl -X PUT "{{uploadImagePresignUrl}}" \
     -H "Content-Type: {{mimeType}}" \
     -H "Content-Length: {{fileSize}}" \
     --data-binary "@test_pic_1.jpg"
```

### Step 4: Create 3 Posts

```bash
for i in {1..3}
do
  curl -X POST "http://0.0.0.0:8010/v1/api/posts" \
       -H "Content-Type: application/json" \
       -H "userId: 507f1f77bcf86cd799439011" \
       -d '{
             "imageFilePath": "{{uploadImagePath from step 3}}",
             "content": "Test post content '"$i"', TestIntegration"
           }'
done
```

### Step 5: Fetch Posts in Default Order

```bash
curl -X GET "http://0.0.0.0:8010/v1/api/posts?limit=50&orderBy=post_id" \
     -H "userId: 507f1f77bcf86cd799439011"
```

### Step 6: Comment on Post 1

```bash
# Assume the first post ID is stored in a variable POST_ID
curl -X POST "http://0.0.0.0:8010/v1/api/posts/{{createdPostID[0]}}/comments" \
     -H "Content-Type: application/json" \
     -H "userId: 507f1f77bcf86cd799439012" \
     -d '{
           "content": "Test comment1 content1 lloremipsum",
           "postId": "{{createdPostID[0]}}"
         }'
```

### Step 7: Fetch Posts in Comment Count Order

```bash
curl -X GET "http://0.0.0.0:8010/v1/api/posts?limit=50&orderBy=comment_count" \
     -H "userId: 507f1f77bcf86cd799439011"
```