package dto

// OrderByType is an enum for order by options.
type OrderByType string

const (
	OrderByCommentCount OrderByType = "comment_count"
	OrderByPostID       OrderByType = "post_id"
	// OrderByCreatedTime   OrderByType = "created_time" // TODO: support other order by options
)

// FetchPostsReq is a struct that represents the FetchPostsReq DTO.
type FetchPostsReq struct {
	Limit          int64       `query:"limit" vd:"$>0"` // Default: 10
	PreviousCursor string      `query:"previousCursor"` // Default: empty
	UserId         string      `header:"userId" vd:"len($)>0"`
	OrderBy        OrderByType `query:"orderBy"` // Default: OrderByPostID
}

type FetchPostsResp struct {
	Posts          []*Post `json:"posts"`          // Fetched posts
	PreviousCursor string  `json:"previousCursor"` // Previous cursor
	HasMore        bool    `json:"hasMore"`        // Whether there are more posts to fetch
}

type Post struct {
	Id                     string     `json:"id"`                // Post ID
	CreatedAtMilli         int64      `json:"createdAt"`         // Post created time in milliseconds
	Content                string     `json:"content"`           // Post content
	CommentCount           int        `json:"commentCount"`      // Post comment count
	RecentComments         []*Comment `json:"recentComments"`    // Recent comments
	RecentCommentedAtMilli int64      `json:"recentCommentedAt"` // Recent commented time in milliseconds
	CreatorId              string     `json:"creatorId"`         // Post creator ID
	CreatorName            string     `json:"creatorName"`       // Post creator name
	ImageId                string     `json:"imageId"`           // Post image ID
	ImageURL               string     `json:"imageURL"`          // Post image URL
	CommentCountCursor     string     `json:"commentCountCursor"`
}

type Comment struct {
	Id             string `json:"id"`          // Comment ID
	CreatedAtMilli int64  `json:"createdAt"`   // Comment created time in milliseconds
	Content        string `json:"content"`     // Comment content
	CreatorId      string `json:"creatorId"`   // Comment creator ID
	CreatorName    string `json:"creatorName"` // Comment creator name (Not Subject to Dynamic Updates)
}
