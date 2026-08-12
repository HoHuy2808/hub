package post

import (
	"time"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	UserId      uuid.UUID `json:"-"`
	Content     *string   `json:"content"`
	Attachments []string  `json:"attachments"`
	IsPublic    *bool     `json:"is_public"`
}

type CreatePostResponse struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Content     *string   `json:"content"`
	Attachments []string  `json:"attachments"`
	IsPublic    *bool     `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

type UpdatePostRequest struct {
	Id       uuid.UUID `json:"-"`
	Content  *string   `json:"content"`
	IsPublic *bool     `json:"is_public"`
}

type UpdatePostResponse struct {
	Id        uuid.UUID `json:"id"`
	Content   *string   `json:"content"`
	IsPublic  *bool     `json:"is_public"`
	UpdatedAt time.Time `json:"updated_at"`
}
