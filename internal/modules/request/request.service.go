package request

import (
	"context"
	"errors"
	"hub/internal/database/entities"
	"hub/internal/modules/contact"
	"hub/internal/modules/notification"
	"hub/pkg/transaction"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequestService interface {
	GetAll(ctx context.Context) ([]RequestResponse, error)
	SendRequest(ctx context.Context, req *SendRequestReq) (*SendRequestRes, error)
	UpdateRequest(ctx context.Context, req *UpdateRequestReq) (*SendRequestRes, error)
	DeleteRequest(ctx context.Context, requestId, userId uuid.UUID) error
}

type RequestServiceImp struct {
	requestRepo      RequestRepository
	notificationRepo notification.NotificationRepository
	contactRepo      contact.ContactRepository
	txManager        transaction.TransactionManager
}

func NewRequestServiceImp(requestRepo RequestRepository, notificationRepo notification.NotificationRepository, contactRepo contact.ContactRepository, txManager transaction.TransactionManager) *RequestServiceImp {
	return &RequestServiceImp{
		requestRepo:      requestRepo,
		notificationRepo: notificationRepo,
		contactRepo:      contactRepo,
		txManager:        txManager,
	}
}

func (r *RequestServiceImp) GetAll(ctx context.Context) ([]RequestResponse, error) {
	reqs, err := r.requestRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var res []RequestResponse
	for _, req := range reqs {
		res = append(res, RequestResponse{
			Id:         req.Id,
			SenderId:   req.SenderId,
			ReceiverId: req.ReceiverId,
			Status:     string(req.Status),
			CreatedAt:  req.CreatedAt,
		})
	}
	return res, nil
}

func (r *RequestServiceImp) SendRequest(ctx context.Context, req *SendRequestReq) (*SendRequestRes, error) {
	newReq := &entities.Request{
		SenderId:   req.SenderId,
		ReceiverId: req.ReceiverId,
		Status:     entities.PENDING,
	}

	if req.SenderId == req.ReceiverId {
		return nil, errors.New("SenderId and ReceiverId cannot be the same")
	}
	_, err := r.contactRepo.FindContact(ctx, req.SenderId, req.ReceiverId)
	if err == nil {
		return nil, errors.New("Already friend")
	}
	alreadyRequest, err := r.requestRepo.FindRequest(ctx, req.SenderId, req.ReceiverId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if alreadyRequest != nil {
		return nil, errors.New("Request already exists")
	}

	err = r.requestRepo.CreateRequest(ctx, newReq)

	if err != nil {
		return nil, err
	}

	notification := &entities.Notification{
		FromUserId: req.SenderId,
		ToUserId:   req.ReceiverId,
		Type:       "request",
		Metadata:   []byte(`{"request_id": "` + newReq.Id.String() + `"}`),
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

	return &SendRequestRes{
		Id:         newReq.Id,
		SenderId:   newReq.SenderId,
		ReceiverId: newReq.ReceiverId,
		Status:     string(newReq.Status),
		CreatedAt:  newReq.CreatedAt,
	}, nil
}

func (r *RequestServiceImp) UpdateRequest(ctx context.Context, req *UpdateRequestReq) (*SendRequestRes, error) {
	existingRequest, err := r.requestRepo.FindRequestById(ctx, req.RequestId)
	if err != nil {
		return nil, errors.New("Request not found")
	}

	if existingRequest.ReceiverId != req.ReceiverId {
		return nil, errors.New("You are not authorized to update this request")
	}

	if existingRequest.Status != entities.PENDING {
		return nil, errors.New("Request is already processed")
	}

	if req.Status != string(entities.ACCEPT) && req.Status != string(entities.REJECT) {
		return nil, errors.New("Invalid status")
	}

	existingRequest.Status = entities.RequestStatus(req.Status)

	err = r.txManager.ExecTx(ctx, func(txCtx context.Context) error {
		if err := r.requestRepo.UpdateRequest(txCtx, existingRequest); err != nil {
			return err
		}

		if existingRequest.Status == entities.ACCEPT {
			err := r.contactRepo.CreateContact(txCtx, &entities.Contact{
				UserId:    existingRequest.SenderId,
				ContactId: existingRequest.ReceiverId,
			})
			if err != nil {
				return err
			}
			err = r.contactRepo.CreateContact(txCtx, &entities.Contact{
				UserId:    existingRequest.ReceiverId,
				ContactId: existingRequest.SenderId,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if existingRequest.Status == entities.ACCEPT {
		notification := &entities.Notification{
			FromUserId: existingRequest.ReceiverId,
			ToUserId:   existingRequest.SenderId,
			Type:       "accept_request",
			Metadata:   []byte(`{"request_id": "` + existingRequest.Id.String() + `"}`),
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
	}

	return &SendRequestRes{
		Id:         existingRequest.Id,
		SenderId:   existingRequest.SenderId,
		ReceiverId: existingRequest.ReceiverId,
		Status:     string(existingRequest.Status),
		CreatedAt:  existingRequest.CreatedAt,
	}, nil
}

func (r *RequestServiceImp) DeleteRequest(ctx context.Context, requestId, userId uuid.UUID) error {
	existingRequest, err := r.requestRepo.FindRequestById(ctx, requestId)
	if err != nil {
		return errors.New("Request not found")
	}
	if existingRequest.SenderId != userId && existingRequest.ReceiverId != userId {
		return errors.New("You are not authorized to delete this request")
	}
	if existingRequest.Status != entities.PENDING {
		return errors.New("Request is already processed")
	}
	err = r.requestRepo.DeleteRequest(ctx, requestId)
	if err != nil {
		return err
	}
	return nil
}
