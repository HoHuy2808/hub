package post

import (
	"context"
	"hub/internal/database/entities"

	"github.com/google/uuid"
)

type PostService interface {
	GetAll(ctx context.Context, params QueryParams) ([]entities.Post, int64, error)
	CreatePost(ctx context.Context, post *CreatePostRequest) (*CreatePostResponse, error)
	UpdatePost(ctx context.Context, post *UpdatePostRequest) (*UpdatePostResponse, error)
	DeletePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error)
}

type PostServiceImp struct {
	postRepo PostRepository
}

func NewPostServiceImp(postRepo PostRepository) *PostServiceImp {
	return &PostServiceImp{postRepo: postRepo}
}
func (p *PostServiceImp) GetAll(ctx context.Context, params QueryParams) ([]entities.Post, int64, error) {
	return p.postRepo.GetAll(ctx, params)
}

func (p *PostServiceImp) CreatePost(ctx context.Context, post *CreatePostRequest) (*CreatePostResponse, error) {
	var listAttachments []entities.PostAttachment
	for _, attachment := range post.Attachments {
		listAttachments = append(listAttachments, entities.PostAttachment{
			MediaUrl: &attachment,
		})
	}

	newPost := &entities.Post{
		UserId:      post.UserId,
		Content:     post.Content,
		Attachments: listAttachments,
		IsPublic:    post.IsPublic,
	}

	if err := p.postRepo.CreatePost(ctx, newPost); err != nil {
		return nil, err
	}
	return &CreatePostResponse{
		Id:          newPost.Id,
		UserId:      newPost.UserId,
		Content:     newPost.Content,
		Attachments: post.Attachments,
		IsPublic:    newPost.IsPublic,
		CreatedAt:   newPost.CreatedAt,
	}, nil
}

func (p *PostServiceImp) UpdatePost(ctx context.Context, post *UpdatePostRequest) (*UpdatePostResponse, error) {
	existingPost, err := p.postRepo.FindPostById(ctx, post.Id)
	if err != nil {
		return nil, err
	}

	if post.Content != nil {
		existingPost.Content = post.Content
	}
	if post.IsPublic != nil {
		existingPost.IsPublic = post.IsPublic
	}

	err = p.postRepo.UpdatePost(ctx, existingPost)
	if err != nil {
		return nil, err
	}
	return &UpdatePostResponse{
		Id:        existingPost.Id,
		Content:   existingPost.Content,
		IsPublic:  existingPost.IsPublic,
		UpdatedAt: existingPost.UpdatedAt,
	}, nil
}

func (p *PostServiceImp) DeletePost(ctx context.Context, postId uuid.UUID, userId uuid.UUID) (bool, error) {
	return p.postRepo.DeletePost(ctx, postId, userId)
}
