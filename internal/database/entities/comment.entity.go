package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	PostId    uuid.UUID `gorm:"type:uuid;not null;"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;"`
	ParentId  *uuid.UUID
	Content   string    `gorm:"type:text;not null;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time
}
