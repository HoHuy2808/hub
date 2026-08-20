package notification

import (
	"hub/internal/database/entities"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notification *entities.Notification) error
}

type NotificationRepositoryImp struct {
	db *gorm.DB
}

func NewNotificationRepositoryImp(db *gorm.DB) *NotificationRepositoryImp {
	return &NotificationRepositoryImp{db: db}
}

func (n *NotificationRepositoryImp) CreateNotification(notification *entities.Notification) error {
	result := n.db.Create(notification)
	return result.Error
}
