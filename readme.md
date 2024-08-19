# Bandlab Assignment

## UseCase Analysis

1. Create Post, with Image Uploading
2. Create comments binded to a post
3. Delete own comments
4. List posts order by comments count

## Resources estimation

1. Storage for images (3MB for original image avg size, 100kb after resize): 1000 *3MB* 24h *365d + 1000* 100KB *24h* 365d ->> **27TB per year**
2. Storage for posts and comments (1KB per post, 0.5KB per comment): 1000 *1KB* 24h *365d + 100K* 1KB *24h* 365d ->> **430GB per year**

## Requirements Analysis and Assumptions

1. Image uploading is a heavy operation, and the requirement wants to support up to **100MB**, with **weak client network**. I strongly not suggest this large size. However, to support this feature we wont process uploading in the server. We will let client to upload image to OSS directly.
2.

## Tech Stack

1. Server: Golang
2. Image Storage: AWS S3 (Use Cloudflare R2 to mock)
3. Database: MongoDB (Lock Container Mock)
4. Image Processing: AWS Lambda (Use a server goroutine to mock)

## Architecture

![Reading](readme/architecutre.png#center)  

## Creating Post

To optimise, I separated the image uploading process from the post creation process. The client will upload the image directly to the S3 bucket, by calling presigned url.  The client will then send the image URL to the server to create the post. The server will then create the post and store the image URL in the database.

## Future Improvements

1. Use multipart upload for weak network / Let client to resize the image down to 5MB before upload
2. Start Image processing when uploading finishes with callback
3. Use a queue to update post comment count + latest comments
