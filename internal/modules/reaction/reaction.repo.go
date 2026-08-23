package reaction

import (
	"context"
	"hub/internal/database/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReactionRepository interface {
	GetAll(ctx context.Context, params *QueryParams) ([]*entities.Reaction, int64, error)
	React(ctx context.Context, reaction *entities.Reaction) error
	FindReaction(ctx context.Context, postId, reactorId uuid.UUID) (*entities.Reaction, error)
	UpdateReact(ctx context.Context, reaction *entities.Reaction) error
	Unreact(ctx context.Context, PostId, ReactorId uuid.UUID) error
}

type ReactionRepositoryImp struct {
	db *gorm.DB
}

func NewReactionRepositoryImp(db *gorm.DB) *ReactionRepositoryImp {
	return &ReactionRepositoryImp{db: db}
}

// Kiểm tra xem trong Context có tx không. Có thì dùng tx, không có thì dùng db gốc.
func (r *ReactionRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("TxKey").(*gorm.DB)
	if ok {
		return tx // Đang nằm trong Transaction
	}
	return r.db // Chạy đơn lẻ bình thường
}

func (r *ReactionRepositoryImp) GetAll(ctx context.Context, params *QueryParams) ([]*entities.Reaction, int64, error) {
	db := r.getDB(ctx)
	var reactions []*entities.Reaction
	var total int64

	query := db.Model(&entities.Reaction{}).Where("post_id = ?", params.PostId)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if params.Limit != 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset != 0 {
		query = query.Offset(params.Offset)
	}

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	result := query.Find(&reactions)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return reactions, total, nil
}

func (r *ReactionRepositoryImp) React(ctx context.Context, reaction *entities.Reaction) error {
	db := r.getDB(ctx)
	result := db.Create(reaction)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ReactionRepositoryImp) FindReaction(ctx context.Context, postId, reactorId uuid.UUID) (*entities.Reaction, error) {
	db := r.getDB(ctx)
	var reaction entities.Reaction
	result := db.Where("post_id = ? AND reactor_id = ?", postId, reactorId).First(&reaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &reaction, nil
}

func (r *ReactionRepositoryImp) UpdateReact(ctx context.Context, reaction *entities.Reaction) error {
	db := r.getDB(ctx)
	result := db.Model(reaction).
		Where("post_id = ? AND reactor_id = ?", reaction.PostId, reaction.ReactorId).
		Updates(map[string]interface{}{"type": reaction.Type})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ReactionRepositoryImp) Unreact(ctx context.Context, PostId, ReactorId uuid.UUID) error {
	db := r.getDB(ctx)
	result := db.Where("post_id = ? AND reactor_id = ?", PostId, ReactorId).Delete(&entities.Reaction{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
