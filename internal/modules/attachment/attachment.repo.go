package attachment

import (
	"context"
	"hub/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	GetAll(ctx context.Context, query QueryParams) ([]entities.PostAttachment, int64, error)
	AddAttachment(ctx context.Context, attachments []entities.PostAttachment) error
	DeleteAttachment(ctx context.Context, id uuid.UUID) error
}

type AttachmentRepositoryImp struct {
	db *gorm.DB
}

func NewAttachmentRepositoryImp(db *gorm.DB) *AttachmentRepositoryImp {
	return &AttachmentRepositoryImp{db: db}
}

func (a *AttachmentRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("TxKey").(*gorm.DB)
	if ok {
		return tx
	}
	return a.db
}

func (a *AttachmentRepositoryImp) GetAll(ctx context.Context, params QueryParams) ([]entities.PostAttachment, int64, error) {
	var attachments []entities.PostAttachment
	var total int64

	db := a.getDB(ctx)
	query := db.Model(&entities.PostAttachment{})

	if params.PostId != (uuid.UUID{}) {
		query = query.Where("post_attachments.post_id = ?", params.PostId)
	}

	if params.UserId != (uuid.UUID{}) {
		query = query.Joins("JOIN posts ON posts.id = post_attachments.post_id").
			Where("posts.user_id = ?", params.UserId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offSet := (params.Page - 1) * params.Limit
	query = query.Limit(params.Limit).Offset(offSet)

	if err := query.Find(&attachments).Error; err != nil {
		return nil, 0, err
	}
	return attachments, total, nil
}

func (a *AttachmentRepositoryImp) AddAttachment(ctx context.Context, attachments []entities.PostAttachment) error {
	if len(attachments) == 0 {
		return nil
	}

	// Batch Insert (Tạo nhiều dòng cùng lúc) khi truyền vào một mảng (slice)
	db := a.getDB(ctx)
	result := db.Create(&attachments)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a *AttachmentRepositoryImp) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	db := a.getDB(ctx)
	return db.Where("id = ?", id).Delete(&entities.PostAttachment{}).Error
}
