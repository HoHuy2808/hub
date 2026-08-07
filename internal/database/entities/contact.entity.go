package entities

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId    uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	ContactId uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex"`
	IsBlocked bool       `gorm:"default:false;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time `gorm:"default:null;"`

	User    *User `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	Contact *User `gorm:"foreignKey:ContactId;references:Id;constraint:OnDelete:CASCADE;"`
}
