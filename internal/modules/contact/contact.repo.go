package contact

import (
	"context"
	"hub/internal/database/entities"
	"hub/pkg/transaction"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ContactRepository interface {
	GetAll(ctx context.Context, params QueryParams) ([]entities.Contact, int64, error)
	CreateContact(ctx context.Context, contact *entities.Contact) error
	FindContact(ctx context.Context, userId, contactId uuid.UUID) (*entities.Contact, error)
	UpdateContact(ctx context.Context, userId, contactId uuid.UUID, isBlocked bool) error
	DeleteContact(ctx context.Context, userId, contactId uuid.UUID) error
}

type ContactRepositoryImp struct {
	db *gorm.DB
}

func NewContactRepositoryImp(db *gorm.DB) *ContactRepositoryImp {
	return &ContactRepositoryImp{db: db}
}

func (c *ContactRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return c.db
}

func (c *ContactRepositoryImp) GetAll(ctx context.Context, params QueryParams) ([]entities.Contact, int64, error) {
	db := c.getDB(ctx)
	var contacts []entities.Contact
	var total int64

	// Join bảng users để lấy thông tin user_name và lọc theo user_id
	query := db.Model(&entities.Contact{}).
		Preload("Contact").
		Joins("JOIN users ON users.id = contacts.contact_id").
		Where("contacts.user_id = ?", params.UserId) // Lọc theo ID của người dùng

	// Tìm kiếm theo tên
	if params.UserName != "" {
		query = query.Where("users.user_name ILIKE ?", "%"+params.UserName+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}
	if params.SortBy != "" {
		query = query.Order("contacts." + params.SortBy + " " + params.SortOrder)
	}

	if err := query.Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (c *ContactRepositoryImp) CreateContact(ctx context.Context, contact *entities.Contact) error {
	db := c.getDB(ctx)
	if err := db.Create(contact).Error; err != nil {
		return err
	}
	return nil
}

func (c *ContactRepositoryImp) FindContact(ctx context.Context, userId, contactId uuid.UUID) (*entities.Contact, error) {
	db := c.getDB(ctx)
	var contact entities.Contact
	if err := db.Where("user_id = ? AND contact_id = ?", userId, contactId).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (c *ContactRepositoryImp) UpdateContact(ctx context.Context, userId, contactId uuid.UUID, isBlocked bool) error {
	db := c.getDB(ctx)
	if err := db.Model(&entities.Contact{}).
		Where("user_id = ? AND contact_id = ?", userId, contactId).
		Update("is_blocked", isBlocked).Error; err != nil {
		return err
	}
	return nil
}

func (c *ContactRepositoryImp) DeleteContact(ctx context.Context, userId, contactId uuid.UUID) error {
	db := c.getDB(ctx)
	// Xóa cả 2 chiều: (A -> B) và (B -> A)
	if err := db.Where("(user_id = ? AND contact_id = ?) OR (user_id = ? AND contact_id = ?)", userId, contactId, contactId, userId).Delete(&entities.Contact{}).Error; err != nil {
		return err
	}
	return nil
}
