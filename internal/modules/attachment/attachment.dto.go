package attachment

import (
	"github.com/google/uuid"
)

type AddAttachmentRequest struct {
	PostId uuid.UUID `json:"-"`
	URLs   []string  `json:"urls" binding:"required"`
}

type AddAttachmentResponse struct {
	PostId uuid.UUID `json:"post_id"`
	URLs   []string  `json:"url"`
}

type QueryParams struct {
	Page   int `form:"page,default=1"`
	Limit  int `form:"limit,default=10"`
	PostId uuid.UUID
	UserId uuid.UUID `form:"-"`
}
