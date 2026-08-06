package entities

import (
	"time"

	"github.com/google/uuid"
)

type PostAttachment struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	PostId    uuid.UUID `gorm:"type:uuid;not null;"`
	MediaUrl  *string   `gorm:"type:string"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time
}
