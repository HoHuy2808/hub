package attachment

import (
	"context"
	"hub/internal/database/entities"

	"github.com/google/uuid"
)

type AttachmentService interface {
	GetAll(ctx context.Context, params QueryParams) ([]entities.PostAttachment, int64, error)
	AddAttachment(ctx context.Context, req *AddAttachmentRequest) (*AddAttachmentResponse, error)
	DeleteAttachment(ctx context.Context, id uuid.UUID) error
}

type AttachmentServiceImp struct {
	attachmentRepo AttachmentRepository
}

func NewAttachmentServiceImp(attachmentRepo AttachmentRepository) *AttachmentServiceImp {
	return &AttachmentServiceImp{attachmentRepo: attachmentRepo}
}

func (a *AttachmentServiceImp) GetAll(ctx context.Context, params QueryParams) ([]entities.PostAttachment, int64, error) {
	return a.attachmentRepo.GetAll(ctx, params)
}

func (a *AttachmentServiceImp) AddAttachment(ctx context.Context, req *AddAttachmentRequest) (*AddAttachmentResponse, error) {
	var attachments []entities.PostAttachment

	for _, url := range req.URLs {
		urlCopy := url
		attachments = append(attachments, entities.PostAttachment{
			PostId:   req.PostId,
			MediaUrl: &urlCopy,
		})
	}

	err := a.attachmentRepo.AddAttachment(ctx, attachments)
	if err != nil {
		return nil, err
	}

	return &AddAttachmentResponse{
		PostId: req.PostId,
		URLs:   req.URLs,
	}, nil
}

func (a *AttachmentServiceImp) DeleteAttachment(ctx context.Context, id uuid.UUID) error {
	return a.attachmentRepo.DeleteAttachment(ctx, id)
}
