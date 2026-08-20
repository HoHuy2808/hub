package comment

import (
	"encoding/json"
	"errors"
	"log"

	"hub/internal/database/entities"
	"hub/internal/modules/notification"
	"hub/internal/modules/post"

	"github.com/google/uuid"
)

type CommentService interface {
	GetAll(params QueryParams) ([]entities.Comment, int64, error)
	CreateComment(comment *CreateCommentRequest) (*CreateCommentResponse, error)
	UpdateComment(comment *UpdateCommentRequest) (*UpdateCommentResponse, error)
	DeleteComment(id uuid.UUID, userId uuid.UUID) (bool, error)
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

func (c *CommentServiceImp) GetAll(params QueryParams) ([]entities.Comment, int64, error) {
	return c.commentRepo.GetAll(params)
}

func (c *CommentServiceImp) CreateComment(req *CreateCommentRequest) (*CreateCommentResponse, error) {
	post, err := c.postRepo.FindPostById(req.PostId)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, errors.New("Post not found")
	}
	comment := &entities.Comment{
		PostId:   req.PostId,
		UserId:   req.UserId,
		ParentId: req.ParentId,
		Content:  req.Content,
	}

	err = c.commentRepo.CreateComment(comment)
	if err != nil {
		return nil, err
	}

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

		var targetUserId uuid.UUID
		if req.ParentId != nil {
			parentComment, err := c.commentRepo.FindCommentById(*req.ParentId)
			if err == nil {
				targetUserId = parentComment.UserId
			}
		} else {
			targetUserId = post.UserId
		}

		if targetUserId == uuid.Nil || targetUserId == req.UserId {
			return
		}

		notification := &entities.Notification{
			FromUserId: req.UserId,
			ToUserId:   targetUserId,
			Type:       "comment",
			Metadata:   json.RawMessage(metaData),
		}
		c.notificationRepo.CreateNotification(notification)
	}()

	return &CreateCommentResponse{
		Id:        comment.Id,
		PostId:    comment.PostId,
		UserId:    comment.UserId,
		ParentId:  comment.ParentId,
		Metadata:  metaDataMap,
		CreatedAt: comment.CreatedAt,
	}, nil
}

func (c *CommentServiceImp) UpdateComment(req *UpdateCommentRequest) (*UpdateCommentResponse, error) {
	comment, err := c.commentRepo.FindCommentById(req.Id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("Comment not found")
	}
	if comment.UserId != req.UserId {
		return nil, errors.New("You are not authorized to update this comment")
	}
	comment.Content = req.Content
	err = c.commentRepo.UpdateComment(comment)
	if err != nil {
		return nil, err
	}
	return &UpdateCommentResponse{
		Id:        comment.Id,
		PostId:    comment.PostId,
		UserId:    comment.UserId,
		Content:   comment.Content,
		UpdatedAt: comment.UpdatedAt,
	}, nil
}

func (c *CommentServiceImp) DeleteComment(id uuid.UUID, userId uuid.UUID) (bool, error) {
	return c.commentRepo.DeleteComment(id, userId)
}
