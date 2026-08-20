package post

import (
	"hub/internal/database/entities"

	"github.com/google/uuid"
)

type PostService interface {
	GetAll(params QueryParams) ([]entities.Post, int64, error)
	CreatePost(post *CreatePostRequest) (*CreatePostResponse, error)
	UpdatePost(post *UpdatePostRequest) (*UpdatePostResponse, error)
	DeletePost(postId uuid.UUID, userId uuid.UUID) (bool, error)
}

type PostServiceImp struct {
	postRepo PostRepository
}

func NewPostServiceImp(postRepo PostRepository) *PostServiceImp {
	return &PostServiceImp{postRepo: postRepo}
}
func (p *PostServiceImp) GetAll(params QueryParams) ([]entities.Post, int64, error) {
	return p.postRepo.GetAll(params)
}

func (p *PostServiceImp) CreatePost(post *CreatePostRequest) (*CreatePostResponse, error) {
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

	if err := p.postRepo.CreatePost(newPost); err != nil {
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

func (p *PostServiceImp) UpdatePost(post *UpdatePostRequest) (*UpdatePostResponse, error) {
	existingPost, err := p.postRepo.FindPostById(post.Id)
	if err != nil {
		return nil, err
	}

	if post.Content != nil {
		existingPost.Content = post.Content
	}
	if post.IsPublic != nil {
		existingPost.IsPublic = post.IsPublic
	}

	err = p.postRepo.UpdatePost(existingPost)
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

func (p *PostServiceImp) DeletePost(postId uuid.UUID, userId uuid.UUID) (bool, error) {
	return p.postRepo.DeletePost(postId, userId)
}
