package entities

import (
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	ACCEPT  RequestStatus = "ACCEPT"
	PENDING RequestStatus = "PENDING"
	REJECT  RequestStatus = "REJECT"
)

type Request struct {
	Id         uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SenderId   uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex"`
	ReceiverId uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex"`
	Status     RequestStatus `gorm:"type:string;not null;default:PENDING;"`
	CreatedAt  time.Time     `gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt  *time.Time    `gorm:"default:null;"`

	Sender   *User `gorm:"foreignKey:SenderId;references:Id;constraint:OnDelete:CASCADE;"`
	Receiver *User `gorm:"foreignKey:ReceiverId;references:Id;constraint:OnDelete:CASCADE;"`
}
