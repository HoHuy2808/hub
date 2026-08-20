package comment

import (
	"hub/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository interface {
	GetAll(params QueryParams) ([]entities.Comment, int64, error)
	CreateComment(comment *entities.Comment) error
	UpdateComment(comment *entities.Comment) error
	DeleteComment(id uuid.UUID, userId uuid.UUID) (bool, error)
	FindCommentById(id uuid.UUID) (*entities.Comment, error)
}

type CommentRepositoryImp struct {
	db *gorm.DB
}

func NewCommentRepositoryImp(db *gorm.DB) *CommentRepositoryImp {
	return &CommentRepositoryImp{db: db}
}

func (c *CommentRepositoryImp) GetAll(params QueryParams) ([]entities.Comment, int64, error) {
	var comments []entities.Comment
	var total int64

	// Chỉ lấy comment của bài viết hiện tại
	query := c.db.Model(&entities.Comment{}).Where("post_id = ?", params.PostId)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Gán giá trị mặc định cho phân trang và sắp xếp
	limit := params.Limit
	if limit <= 0 {
		limit = 10 // Mặc định trả về 10 comment mỗi trang
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := params.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}

	orderBy := sortBy + " " + sortOrder
	query = query.Order(orderBy).Limit(limit).Offset(params.Offset)

	if err := query.Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func (c *CommentRepositoryImp) CreateComment(comment *entities.Comment) error {
	result := c.db.Create(comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *CommentRepositoryImp) UpdateComment(comment *entities.Comment) error {
	result := c.db.Model(comment).Updates(map[string]interface{}{
		"content": comment.Content,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (c *CommentRepositoryImp) DeleteComment(id uuid.UUID, userId uuid.UUID) (bool, error) {
	result := c.db.Where("id = ? AND user_id = ?", id, userId).Delete(&entities.Comment{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (c *CommentRepositoryImp) FindCommentById(id uuid.UUID) (*entities.Comment, error) {
	var comment entities.Comment
	result := c.db.First(&comment, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}
