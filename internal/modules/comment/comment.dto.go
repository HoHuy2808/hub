package comment

import (
	"time"

	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	PostId      uuid.UUID  `json:"-"`
	CommenterId uuid.UUID  `json:"-"`
	ParentId    *uuid.UUID `json:"parent_id"`
	Content     string     `json:"content"`
}

type CreateCommentResponse struct {
	Id          uuid.UUID              `json:"id"`
	PostId      uuid.UUID              `json:"post_id"`
	PostOwnerId uuid.UUID              `json:"post_owner_id"`
	CommenterId uuid.UUID              `json:"user_id"`
	ParentId    *uuid.UUID             `json:"parent_id"`
	Metadata    map[string]interface{} `json:"meta_data"`
	CreatedAt   time.Time              `json:"created_at"`
}

type UpdateCommentRequest struct {
	Id          uuid.UUID `json:"-"`
	CommenterId uuid.UUID `json:"-"`
	Content     string    `json:"content"`
}

type UpdateCommentResponse struct {
	Id          uuid.UUID `json:"id"`
	PostId      uuid.UUID `json:"post_id"`
	CommenterId uuid.UUID `json:"user_id"`
	Content     string    `json:"content"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueryParams struct {
	PostId    uuid.UUID `json:"-"`
	Limit     int       `form:"limit"`
	Offset    int       `form:"off_set"`
	SortBy    string    `form:"sort_by,default=created_at"`
	SortOrder string    `form:"sort_order,default=desc"`
}
