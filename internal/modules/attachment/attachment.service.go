package attachment

import (
	"hub/internal/database/entities"

	"github.com/google/uuid"
)

type AttachmentService interface {
	GetAll(params QueryParams) ([]entities.PostAttachment, int64, error)
	AddAttachment(req *AddAttachmentRequest) (*AddAttachmentResponse, error)
	DeleteAttachment(id uuid.UUID) error
}

type AttachmentServiceImp struct {
	attachmentRepo AttachmentRepository
}

func NewAttachmentServiceImp(attachmentRepo AttachmentRepository) *AttachmentServiceImp {
	return &AttachmentServiceImp{attachmentRepo: attachmentRepo}
}

func (a *AttachmentServiceImp) GetAll(params QueryParams) ([]entities.PostAttachment, int64, error) {
	return a.attachmentRepo.GetAll(params)
}

func (a *AttachmentServiceImp) AddAttachment(req *AddAttachmentRequest) (*AddAttachmentResponse, error) {
	var attachments []entities.PostAttachment

	for _, url := range req.URLs {
		urlCopy := url
		attachments = append(attachments, entities.PostAttachment{
			PostId:   req.PostId,
			MediaUrl: &urlCopy,
		})
	}

	err := a.attachmentRepo.AddAttachment(attachments)
	if err != nil {
		return nil, err
	}

	return &AddAttachmentResponse{
		PostId: req.PostId,
		URLs:   req.URLs,
	}, nil
}

func (a *AttachmentServiceImp) DeleteAttachment(id uuid.UUID) error {
	return a.attachmentRepo.DeleteAttachment(id)
}
