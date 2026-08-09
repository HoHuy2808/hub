package entities

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`

	UserId uuid.UUID

	RefreshToken string `gorm:"type:text;not null;"`

	ExpiresAt time.Time `gorm:"type:timestamp;not null;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
}
