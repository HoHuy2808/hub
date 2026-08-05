package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	UserId      uuid.UUID
	FullName    string
	DateOfBirth *time.Time
	Gender      *string
	AvatarUrl   *string   `gorm:"type:varchar(255);"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
	UpdatedAt   time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
}
