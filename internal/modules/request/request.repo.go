package request

import (
	"context"
	"hub/internal/database/entities"
	"hub/pkg/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestRepository interface {
	GetAll(ctx context.Context) ([]entities.Request, error)
	CreateRequest(ctx context.Context, req *entities.Request) error
	FindRequest(ctx context.Context, senderId, receiverId uuid.UUID) (*entities.Request, error)
	FindRequestById(ctx context.Context, id uuid.UUID) (*entities.Request, error)
	UpdateRequest(ctx context.Context, req *entities.Request) error
	DeleteRequest(ctx context.Context, requestId uuid.UUID) error
}

type RequestRepositoryImp struct {
	db *gorm.DB
}

func NewRequestRepositoryImp(db *gorm.DB) *RequestRepositoryImp {
	return &RequestRepositoryImp{db: db}
}

func (r *RequestRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return r.db
}

func (r *RequestRepositoryImp) GetAll(ctx context.Context) ([]entities.Request, error) {
	db := r.getDB(ctx)
	var reqs []entities.Request
	if err := db.Find(&reqs).Error; err != nil {
		return nil, err
	}
	return reqs, nil
}

func (r *RequestRepositoryImp) CreateRequest(ctx context.Context, req *entities.Request) error {
	db := r.getDB(ctx)
	if err := db.Create(req).Error; err != nil {
		return err
	}
	return nil
}

func (r *RequestRepositoryImp) FindRequest(ctx context.Context, senderId, receiverId uuid.UUID) (*entities.Request, error) {
	db := r.getDB(ctx)
	var req entities.Request
	if err := db.Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *RequestRepositoryImp) FindRequestById(ctx context.Context, id uuid.UUID) (*entities.Request, error) {
	db := r.getDB(ctx)
	var req entities.Request
	if err := db.Where("id = ?", id).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *RequestRepositoryImp) UpdateRequest(ctx context.Context, req *entities.Request) error {
	db := r.getDB(ctx)
	if err := db.Model(&entities.Request{}).
		Where("id = ?", req.Id).
		Updates(map[string]interface{}{
			"status": req.Status,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *RequestRepositoryImp) DeleteRequest(ctx context.Context, requestId uuid.UUID) error {
	db := r.getDB(ctx)
	if err := db.Where("id = ?", requestId).Delete(&entities.Request{}).Error; err != nil {
		return err
	}
	return nil
}
