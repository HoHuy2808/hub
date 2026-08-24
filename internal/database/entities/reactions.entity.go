package entities

import (
	"time"

	"github.com/google/uuid"
)

type ReactionsType string

const (
	Like  ReactionsType = "like"
	Love  ReactionsType = "love"
	Haha  ReactionsType = "haha"
	Wow   ReactionsType = "wow"
	Sad   ReactionsType = "sad"
	Angry ReactionsType = "angry"
)

type Reaction struct {
	Id        uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid();"`
	PostId    uuid.UUID     `gorm:"type:uuid;not null;"`
	ReactorId uuid.UUID     `gorm:"type:uuid;not null;"`
	Type      ReactionsType `gorm:"type:string;not null;"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP;"`
}
