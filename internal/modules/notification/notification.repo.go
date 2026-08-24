package notification

import (
	"context"
	"hub/internal/database/entities"
	"hub/pkg/transaction"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *entities.Notification) error
}

type NotificationRepositoryImp struct {
	db *gorm.DB
}

func NewNotificationRepositoryImp(db *gorm.DB) *NotificationRepositoryImp {
	return &NotificationRepositoryImp{db: db}
}

func (n *NotificationRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return n.db
}

func (n *NotificationRepositoryImp) CreateNotification(ctx context.Context, notification *entities.Notification) error {
	db := n.getDB(ctx)
	result := db.Create(notification)
	return result.Error
}
