package dto

// OrderByType is an enum for order by options.
type OrderByType string

const (
	OrderByCreatedTime  OrderByType = "created_time"
	OrderByCommentCount OrderByType = "comment_count"
)

// FetchPostsReq is a struct that represents the FetchPostsReq DTO.
type FetchPostsReq struct {
	Limit          int         `query:"limit" vd:"$>0"`
	PreviousCursor string      `query:"previousCursor"`
	UserId         string      `header:"userId" vd:"len($)>0"`
	OrderBy        OrderByType `query:"orderBy"`
}

type FetchPostsResp struct {
	Posts          []*Post `json:"posts"`
	HasMore        bool    `json:"hasMore"`
	NextCursor     string  `json:"nextCursor"`
	PreviousCursor string  `json:"previousCursor"`
}

type Post struct {
	Id                     string     `json:"id"`
	CreatedAtMilli         int64      `json:"createdAt"`
	Content                string     `json:"content"`
	CommentCount           int        `json:"commentCount"`
	RecentComments         []*Comment `json:"recentComments"`
	RecentCommentedAtMilli int64      `json:"recentCommentedAt"`
	CreatorId              string     `json:"creatorId"`
	CreatorName            string     `json:"creatorName"`
	ImageId                string     `json:"imageId"`
	ImageURL               string     `json:"imageURL"`
}

type Comment struct {
	Id             string `json:"id"`
	CreatedAtMilli int64  `json:"createdAt"`
	Content        string `json:"content"`
	CreatorId      string `json:"creatorId"`
	CreatorName    string `json:"creatorName"`
}
