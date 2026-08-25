package contact

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ContactService interface {
	GetAll(ctx context.Context, params QueryParams) ([]ContactResponse, int64, error)
	UpdateContact(ctx context.Context, req *UpdateContact) (*UpdateContactRes, error)
	DeleteContact(ctx context.Context, userId, contactId uuid.UUID) error
}

type ContactServiceImp struct {
	contactRepo ContactRepository
}

func NewContactServiceImp(contactRepo ContactRepository) *ContactServiceImp {
	return &ContactServiceImp{contactRepo: contactRepo}
}

func (c *ContactServiceImp) GetAll(ctx context.Context, params QueryParams) ([]ContactResponse, int64, error) {
	contacts, total, err := c.contactRepo.GetAll(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var res []ContactResponse
	for _, contact := range contacts {
		userName := ""
		if contact.Contact != nil {
			userName = contact.Contact.UserName
		}
		res = append(res, ContactResponse{
			Id:        contact.Id,
			ContactId: contact.ContactId,
			UserName:  userName,
			IsBlocked: contact.IsBlocked,
			CreatedAt: contact.CreatedAt,
		})
	}
	return res, total, nil
}

func (c *ContactServiceImp) UpdateContact(ctx context.Context, req *UpdateContact) (*UpdateContactRes, error) {
	err := c.contactRepo.UpdateContact(ctx, req.UserId, req.ContactId, req.IsBlocked)
	if err != nil {
		return nil, err
	}
	return &UpdateContactRes{
		Id:        req.Id,
		ContactId: req.ContactId,
		IsBlocked: req.IsBlocked,
		UpdatedAt: time.Now(),
	}, nil
}

func (c *ContactServiceImp) DeleteContact(ctx context.Context, userId, contactId uuid.UUID) error {
	return c.contactRepo.DeleteContact(ctx, userId, contactId)
}
