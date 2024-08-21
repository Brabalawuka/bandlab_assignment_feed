# Bandlab Assignment

## UseCase Analysis

1. Create Post, with Image Uploading
2. Create comments binded to a post
3. Delete own comments
4. List posts order by comments count

## Resources estimation 5 years

1. Storage for images (Original: 3MB avg size, Resized: 100KB avg size): 1000 *3MB* 24h *365d + 1000* 100KB *24h* 365d * 5y ->> **135TB**
2. Storage for posts (1KB per post): 1000 *1KB* 24h *365d* 5year ->> **43GB** No parition is required
3. Storage for comments (0.5KB per comment): 100K *0.5KB* 24h *365d* 5year ->> **2088GB** (Parition by POST Id may be needed, may cause hot partition issue)
4. Queries - Query Post: 100RPS - Create Post: <1 RPS - Create Comment: 100K / 3600s ->> 28 RPS
   CPU Usage: Assume 50ms per query, 50% time in I/O, 30% CPU usage -> **3 * 4CPU PODs**

## Requirements Pitfalls and Assumptions

1. Image uploading is heavy operation, up to **100MB** single imag with **weak client network**. I strongly not suggest this large size. However, to support this feature we wont uploading image to the server. We will use **presigned URL** for direct client upload to OSS, then process the image in the background.
2. Posts retrival are in desc order of **Comment Count**. However, post may have duplicated comment count and it's hard to sort based on cursor approach. Thus, we will design a **composite sorting** cursor (comment count, last comment time, part of post Id).
3. Updating comment count of posts in parallel may cause issue, thus we will use **optimistic lock** to update the post for MVP. However, this may cause issue in the future if the server is under heavy load.

## Highlights

1. **Composite Sorting Cursor** is used to sort the posts by comment count, last comment time and part of creatorID.
2. **Optimistic Lock** is used to update the post for MVP.
3. **Organised Code** Stucture **handler--sevice--dal** and fulfill single responsibility principle.
4. **Error Handling** is implemented using **custom error** with **error code** and **error message**.
5. **Pre-signed URL** is used to upload image to OSS.

## Tech Stack

1. Server: Golang
2. Image Storage: AWS S3 (Use Cloudflare R2 to mock)
3. Database: MongoDB (User a local Container to Mock) MongoDB is used as it supports dynamic schema and is easy to scale.
4. Image Processing: AWS Lambda (Use a server goroutine to mock)

## Architecture

![Reading](readme/architecutre.png#center)  

## ER Diagram

- A composite sorting key is used to sort the posts as the cursor. The key is a combination of comment count, last comment time and part of creatorID. This ensures no duplicated sorting key.
- Composite Key Structure 12 Bytes:
  - 1-4 bytes: commentCount (uint32, 3.1Billion comments)
  - 5-8 bytes: lastCommentAt (unix timestamp, till 2109)
  - 9-16 bytes: first 8 bytes of PostID 
 Encode as Base64 string
- Storage of image using **path** instead of **url** as the URL should be dynamically generated due to possible change of CDN domain name.

![ER](readme/er.jpeg#center)

## Creating Post

To optimise, I separated the image uploading process from the post creation process. The client will upload the image directly to the S3 bucket, by calling presigned url. The client will then send the allocated image Id to the server to create the post. The server will then create the post and store the image path in the database. Server will then trigger a lambda function to process the image. The lambda function will resize the image and store it back to the S3 bucket. The lambda will then update the post with the resized image URL.

![POST](readme/create_post.jpeg#center)

- **Presigned URL** is generated with specific expire time, content type and length.
- **Delayed Post Status** After the creation of post, other people will not be able to see post image until the image is processed.

## Creating Comment

When a comment is created, the server will not only save the comment to DB, it will async update the post with the latest comment, comment count and the composit sorting key (cursor) as well as the last comment time. This is implemented using a optimistic lock to do the update. If the update fails, the server will retry the update, however, this may cause issue in the future if the server is under heavy load. It may also have consistency issue. A queue may be needed to handle this.

![COMMENT](readme/create_comment.jpeg#center)

## Deleting Comments

Delete comment will mark the comment deleted, meanwhile async modify comment count and latest comment.

![COMMENT](readme/delete_comment.jpeg#center)

## Fetching Posts

To fetch the posts, the server will use the composite sorting cursor to fetch the posts. The server will fetch the posts in batches and return the posts to the client. The client will then use the last post in the batch to fetch the next batch of posts. Since the sorting key updates dynamically, there is possibility of duplicated posts and missing posts. 

![FETCH](readme/fetch.jpeg#center)

- Client will fetch the posts until no more posts are returned.
- Client will have to deduplicate the posts by the post id, since the sorting cursor is updated dynamically.

## Future Improvements

1. Use multipart upload for weak network / Let client to resize the image down to 5MB before upload
2. Start Image processing when uploading finishes with callback
3. Use a queue to update post comment count + latest comments
4. Cache hot posts and comments
5. Aggregate the udpating comment count in the post document to reduce the update frequency
6. Create a collection to store image metadata and delete the image if the post is not posted

## Things to take note before production

1. Replace mock with real services
2. Add CICD and Test coverages and finish all UNITS tests
3. Add monitoring + alerting / logging + tracing for the services and DB
4. Add rate limiting and authentication for the services
5. **Sync with PM to ensure all compromisation is accepted and could be fixed in the future**
