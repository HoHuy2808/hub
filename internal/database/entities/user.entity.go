package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	UserName  string    `gorm:"type:varchar(255);not null;"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex;"`
	Phone     string    `gorm:"type:varchar(255);not null;uniqueIndex;"`
	Password  string    `gorm:"type:varchar(255);not null;"`
	IsActive  bool      `gorm:"type:boolean;default:true;"`
	IsBlock   bool      `gorm:"type:boolean;default:false;"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`

	Profile      UserProfile    `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	RefreshToken []RefreshToken `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	Posts        []Post         `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	Reaction     []Reaction     `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	Comments     []Comment      `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`

	SentRequest     []Request `gorm:"foreignKey:SenderId;references:Id;constraint:OnDelete:CASCADE;"`
	ReceivedRequest []Request `gorm:"foreignKey:ReceiverId;references:Id;constraint:OnDelete:CASCADE;"`

	Contacts        []Contact `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;"`
	ContactsOfUsers []Contact `gorm:"foreignKey:ContactId;references:Id;constraint:OnDelete:CASCADE;"`
}
