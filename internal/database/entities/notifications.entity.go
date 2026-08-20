package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	FromUserId uuid.UUID
	ToUserId   uuid.UUID
	Type       string          `gorm:"type:varchar(255)"`
	Metadata   json.RawMessage `gorm:"type:jsonb"`
	CreatedAt  time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
	// Description string          `gorm:"type:text;"`
	// IsRead     bool `gorm:"type:boolean"`
	// ReadAt     time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;"`
}
