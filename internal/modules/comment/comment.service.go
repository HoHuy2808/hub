package comment

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"hub/internal/database/entities"
	"hub/internal/modules/notification"
	"hub/internal/modules/post"

	"github.com/google/uuid"
)

type CommentService interface {
	GetAll(ctx context.Context, params QueryParams) ([]entities.Comment, int64, error)
	CreateComment(ctx context.Context, comment *CreateCommentRequest) (*CreateCommentResponse, error)
	UpdateComment(ctx context.Context, comment *UpdateCommentRequest) (*UpdateCommentResponse, error)
	DeleteComment(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error)
}

type CommentServiceImp struct {
	commentRepo      CommentRepository
	postRepo         post.PostRepository
	notificationRepo notification.NotificationRepository
}

func NewCommentServiceImp(
	commentRepo CommentRepository,
	postRepo post.PostRepository,
	notificationRepo notification.NotificationRepository,
) *CommentServiceImp {
	return &CommentServiceImp{
		commentRepo:      commentRepo,
		postRepo:         postRepo,
		notificationRepo: notificationRepo,
	}
}

func (c *CommentServiceImp) GetAll(ctx context.Context, params QueryParams) ([]entities.Comment, int64, error) {
	return c.commentRepo.GetAll(ctx, params)
}

func (c *CommentServiceImp) CreateComment(ctx context.Context, req *CreateCommentRequest) (*CreateCommentResponse, error) {
	post, err := c.postRepo.FindPostById(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Post not found")
	}

	// Tạo comment
	comment := &entities.Comment{
		PostId:      req.PostId,
		CommenterId: req.CommenterId,
		ParentId:    req.ParentId,
		Content:     req.Content,
	}

	err = c.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Gửi thông báo đến chủ post
	metaDataMap := make(map[string]interface{})
	metaDataMap["post_id"] = req.PostId
	metaDataMap["comment_id"] = comment.Id
	metaDataMap["content"] = req.Content
	metaData, _ := json.Marshal(metaDataMap)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered in CreateComment notification", r)
			}
		}()

		bgCtx := context.Background()
		var targetUserId uuid.UUID
		if req.ParentId != nil {
			parentComment, err := c.commentRepo.FindCommentById(bgCtx, *req.ParentId)
			if err == nil {
				targetUserId = parentComment.CommenterId
			}
		} else {
			targetUserId = post.UserId
		}

		if targetUserId == uuid.Nil || targetUserId == req.CommenterId {
			return
		}

		notification := &entities.Notification{
			FromUserId: req.CommenterId,
			ToUserId:   targetUserId,
			Type:       "comment",
			Metadata:   json.RawMessage(metaData),
		}

		_ = c.notificationRepo.CreateNotification(bgCtx, notification)
	}()

	return &CreateCommentResponse{
		Id:          comment.Id,
		PostId:      comment.PostId,
		PostOwnerId: post.UserId,
		CommenterId: comment.CommenterId,
		ParentId:    comment.ParentId,
		Metadata:    metaDataMap,
		CreatedAt:   comment.CreatedAt,
	}, nil
}

func (c *CommentServiceImp) UpdateComment(ctx context.Context, req *UpdateCommentRequest) (*UpdateCommentResponse, error) {
	comment, err := c.commentRepo.FindCommentById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("Comment not found")
	}
	if comment.CommenterId != req.CommenterId {
		return nil, errors.New("You are not authorized to update this comment")
	}
	comment.Content = req.Content
	err = c.commentRepo.UpdateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	return &UpdateCommentResponse{
		Id:          comment.Id,
		PostId:      comment.PostId,
		CommenterId: comment.CommenterId,
		Content:     comment.Content,
		UpdatedAt:   comment.UpdatedAt,
	}, nil
}

func (c *CommentServiceImp) DeleteComment(ctx context.Context, id uuid.UUID, userId uuid.UUID) (bool, error) {
	return c.commentRepo.DeleteComment(ctx, id, userId)
}
