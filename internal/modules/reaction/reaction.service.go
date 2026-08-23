package reaction

import (
	"context"
	"errors"
	"hub/internal/database/entities"
	"hub/internal/modules/notification"
	"hub/internal/modules/post"
	"hub/pkg/transaction"
	"log"

	"github.com/google/uuid"
)

type ReactionService interface {
	GetAll(ctx context.Context, params *QueryParams) ([]*entities.Reaction, int64, error)
	React(ctx context.Context, reaction *ReactRequest) (*ReactResponse, error)
	UpdateReact(ctx context.Context, reaction *UpdateReactRequest) (*UpdateReactResponse, error)
	Unreact(ctx context.Context, PostId, ReactorId uuid.UUID) error
}

type ReactionServiceImp struct {
	reactionRepo     ReactionRepository
	postRepo         post.PostRepository
	notificationRepo notification.NotificationRepository
	txManager        transaction.TransactionManager
}

func NewReactionServiceImp(
	reactionRepo ReactionRepository,
	postRepo post.PostRepository,
	notificationRepo notification.NotificationRepository,
	txManager transaction.TransactionManager,
) *ReactionServiceImp {
	return &ReactionServiceImp{
		reactionRepo:     reactionRepo,
		postRepo:         postRepo,
		notificationRepo: notificationRepo,
		txManager:        txManager,
	}
}

func (r *ReactionServiceImp) GetAll(ctx context.Context, params *QueryParams) ([]*entities.Reaction, int64, error) {
	reactions, total, err := r.reactionRepo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return reactions, total, nil
}

func (r *ReactionServiceImp) React(ctx context.Context, req *ReactRequest) (*ReactResponse, error) {

	post, err := r.postRepo.FindPostById(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Post not found")
	}

	_, err = r.reactionRepo.FindReaction(ctx, req.PostId, req.ReactorId)
	if err == nil {
		return nil, errors.New("Reaction already exists")
	}
	reaction := &entities.Reaction{
		PostId:    req.PostId,
		ReactorId: req.ReactorId,
		Type:      entities.ReactionsType(req.Type),
	}

	notification := &entities.Notification{
		FromUserId: req.ReactorId,
		ToUserId:   post.UserId,
		Type:       "react",
		Metadata:   []byte(`{"post_id": "` + req.PostId.String() + `", "reaction_type": "` + req.Type + `"}`),
	}

	err = r.txManager.ExecTx(ctx, func(txCtx context.Context) error {
		// Gọi logic tạo reaction với txCtx thay vì ctx gốc
		if err := r.reactionRepo.React(txCtx, reaction); err != nil {
			return err
		}
		post.TotalReaction += 1
		if err := r.postRepo.UpdatePost(txCtx, post); err != nil {
			return err
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in CreateNotification", r)
				}
			}()
			bgCtx := context.Background()
			if err := r.notificationRepo.CreateNotification(bgCtx, notification); err != nil {
				log.Println("CreateNotification error:", err)
			}
		}()
		return nil
	})

	return &ReactResponse{
		Id:          reaction.Id,
		PostId:      reaction.PostId,
		PostOwnerId: post.UserId,
		ReactorId:   reaction.ReactorId,
		Type:        string(reaction.Type),
		CreatedAt:   reaction.CreatedAt,
	}, nil
}

func (r *ReactionServiceImp) UpdateReact(ctx context.Context, req *UpdateReactRequest) (*UpdateReactResponse, error) {
	post, err := r.postRepo.FindPostById(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Post not found")
	}
	existingReaction, err := r.reactionRepo.FindReaction(ctx, req.PostId, req.ReactorId)
	if err != nil {
		return nil, errors.New("Reaction not found")
	}

	if existingReaction.Type == entities.ReactionsType(req.Type) {
		return nil, errors.New("Reaction type is already the same")
	}

	reaction := &entities.Reaction{
		PostId:    req.PostId,
		ReactorId: req.ReactorId,
		Type:      entities.ReactionsType(req.Type),
	}

	if err := r.reactionRepo.UpdateReact(ctx, reaction); err != nil {
		return nil, err
	}

	return &UpdateReactResponse{
		Id:          reaction.Id,
		PostId:      reaction.PostId,
		PostOwnerId: post.UserId,
		Type:        string(reaction.Type),
		UpdatedAt:   reaction.UpdatedAt,
	}, nil
}

func (r *ReactionServiceImp) Unreact(ctx context.Context, PostId, ReactorId uuid.UUID) error {
	post, err := r.postRepo.FindPostById(ctx, PostId)
	if err != nil {
		return err
	}
	if post == nil {
		return errors.New("Post not found")
	}

	err = r.txManager.ExecTx(ctx, func(txCtx context.Context) error {
		if err := r.reactionRepo.Unreact(txCtx, PostId, ReactorId); err != nil {
			return err
		}
		post.TotalReaction -= 1
		if err := r.postRepo.UpdatePost(txCtx, post); err != nil {
			return err
		}
		return nil
	})

	return err
}
