package post

import (
	"context"
	"hub/internal/database/entities"
	"hub/pkg/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll(ctx context.Context, params QueryParams) ([]entities.Post, int64, error)
	CreatePost(ctx context.Context, post *entities.Post) error
	UpdatePost(ctx context.Context, post *entities.Post) error
	FindPostById(ctx context.Context, postId uuid.UUID) (*entities.Post, error)
	DeletePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error)
}

type PostRepositoryImp struct {
	db *gorm.DB
}

func NewPostRepositoryImp(db *gorm.DB) *PostRepositoryImp {
	return &PostRepositoryImp{db: db}
}

func (p *PostRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transaction.TxKey).(*gorm.DB)
	if ok {
		return tx
	}
	return p.db
}

func (p *PostRepositoryImp) GetAll(ctx context.Context, params QueryParams) ([]entities.Post, int64, error) {
	var posts []entities.Post
	var total int64

	// Tạo câu truy vấn ban đầu
	db := p.getDB(ctx)
	query := db.Model(&entities.Post{})

	// Filter
	if params.Search != "" {
		query = query.Where("content LIKE ?", "%"+params.Search+"%")
	}
	// Đếm tổng trước khi phân trang
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sắp xếp
	orderBy := params.SortBy + " " + params.SortOrder
	query = query.Order(orderBy)

	// Phân trang
	offset := (params.Page - 1) * params.Limit
	query = query.Limit(params.Limit).Offset(offset)

	// Lấy dữ liệu và map vào biến posts
	if err := query.Preload("Attachments").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (p *PostRepositoryImp) CreatePost(ctx context.Context, post *entities.Post) error {
	db := p.getDB(ctx)
	result := db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *PostRepositoryImp) UpdatePost(ctx context.Context, post *entities.Post) error {
	db := p.getDB(ctx)
	result := db.Model(post).Updates(map[string]interface{}{
		"content":        post.Content,
		"is_public":      post.IsPublic,
		"total_reaction": post.TotalReaction,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *PostRepositoryImp) FindPostById(ctx context.Context, postId uuid.UUID) (post *entities.Post, err error) {
	db := p.getDB(ctx)
	result := db.Find(&post, postId)
	if result.Error != nil {
		return nil, result.Error
	}
	return post, nil
}

func (p *PostRepositoryImp) DeletePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error) {
	db := p.getDB(ctx)
	// result := db.Delete(&entities.Post{}, postId)
	result := db.Where("id = ? AND user_id = ?", postId, userId).Delete(&entities.Post{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
