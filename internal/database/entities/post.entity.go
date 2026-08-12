package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	UserId    uuid.UUID `gorm:"type:uuid;not null;"`
	Content   *string   `gorm:"type:text;"`
	IsPublic  *bool     `gorm:"type:boolean;default:true;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time

	Attachments []PostAttachment `gorm:"foreignKey:PostId;references:Id;constraint:OnDelete:CASCADE;"`
	Reactions   []Reaction       `gorm:"foreignKey:PostId;references:Id;constraint:OnDelete:CASCADE;"`
	Comments    []Comment        `gorm:"foreignKey:PostId;references:Id;constraint:OnDelete:CASCADE;"`
}
