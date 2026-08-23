package reaction

import (
	"time"

	"github.com/google/uuid"
)

type ReactRequest struct {
	PostId    uuid.UUID `json:"-"`
	ReactorId uuid.UUID `json:"-"`
	Type      string    `json:"type"`
}

type ReactResponse struct {
	Id          uuid.UUID `json:"id"`
	PostId      uuid.UUID `json:"post_id"`
	PostOwnerId uuid.UUID `json:"post_owner_id"`
	ReactorId   uuid.UUID `json:"-"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdateReactRequest struct {
	PostId    uuid.UUID `json:"-"`
	ReactorId uuid.UUID `json:"-"`
	Type      string    `json:"type"`
}

type UpdateReactResponse struct {
	Id          uuid.UUID `json:"id"`
	PostId      uuid.UUID `json:"post_id"`
	PostOwnerId uuid.UUID `json:"post_owner_id"`
	Type        string    `json:"type"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QueryParams struct {
	PostId uuid.UUID `form:"-"`
	Limit  int       `form:"limit"`
	Offset int       `form:"offset"`
	Type   string    `form:"type"`
}
