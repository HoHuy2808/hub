package request

import (
	"time"

	"github.com/google/uuid"
)

type SendRequestReq struct {
	SenderId   uuid.UUID `json:"-"`
	ReceiverId uuid.UUID `json:"receiver_id"`
}

type SendRequestRes struct {
	Id         uuid.UUID `json:"id"`
	SenderId   uuid.UUID `json:"-"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type UpdateRequestReq struct {
	RequestId  uuid.UUID `json:"-"`
	ReceiverId uuid.UUID `json:"-"`
	Status     string    `json:"status" binding:"required"`
}

type RequestResponse struct {
	Id         uuid.UUID `json:"id"`
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
