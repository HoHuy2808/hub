package contact

import (
	"time"

	"github.com/google/uuid"
)

type UpdateContact struct {
	Id        uuid.UUID `json:"-"`
	UserId    uuid.UUID `json:"-"`
	ContactId uuid.UUID `json:"-"`
	IsBlocked bool      `json:"is_block"`
}

type UpdateContactRes struct {
	Id        uuid.UUID `json:"id"`
	ContactId uuid.UUID `json:"contact_id"`
	IsBlocked bool      `json:"is_block"`
	UpdatedAt time.Time `json:"updated_at"`
}

type QueryParams struct {
	UserId    string `form:"userId"`
	UserName  string `form:"userName"`
	Limit     int    `form:"limit"`
	Offset    int    `form:"offset"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
}

type ContactResponse struct {
	Id        uuid.UUID `json:"id"`
	ContactId uuid.UUID `json:"contact_id"`
	UserName  string    `json:"user_name"`
	IsBlocked bool      `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginatedContactResponse struct {
	Data  []ContactResponse `json:"data"`
	Total int64             `json:"total"`
}
