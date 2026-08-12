package post

import (
	"hub/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(post *entities.Post) error
	UpdatePost(post *entities.Post) error
	FindPostById(postId uuid.UUID) (*entities.Post, error)
	DeletePost(postId uuid.UUID, userId uuid.UUID) (bool, error)
}

type PostRepositoryImp struct {
	db *gorm.DB
}

func NewPostRepositoryImp(db *gorm.DB) *PostRepositoryImp {
	return &PostRepositoryImp{db: db}
}

func (p *PostRepositoryImp) CreatePost(post *entities.Post) error {
	result := p.db.Create(post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *PostRepositoryImp) UpdatePost(post *entities.Post) error {
	// result := p.db.Save(post)
	result := p.db.Model(post).Updates(map[string]interface{}{
		"content":   post.Content,
		"is_public": post.IsPublic,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *PostRepositoryImp) FindPostById(postId uuid.UUID) (post *entities.Post, err error) {
	result := p.db.Find(&post, postId)
	if result.Error != nil {
		return nil, result.Error
	}
	return post, nil
}

func (p *PostRepositoryImp) DeletePost(postId uuid.UUID, userId uuid.UUID) (bool, error) {
	// result := p.db.Delete(&entities.Post{}, postId)
	result := p.db.Where("id = ? AND user_id = ?", postId, userId).Delete(&entities.Post{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
