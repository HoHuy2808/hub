package auth

import (
	"errors"
	"hub/internal/database/entities"

	"gorm.io/gorm"
)

type UserService interface {
	Register(req *RegisterRequest) (*RegisterResponse, error)
}

type UserServiceImp struct {
	userRepo UserRepository
}

func NewUserServiceImp(repo UserRepository) *UserServiceImp {
	return &UserServiceImp{userRepo: repo}
}

func (u *UserServiceImp) Register(req *RegisterRequest) (*RegisterResponse, error) {
	emailUser, err := u.userRepo.FindUserByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if emailUser != nil {
		return nil, errors.New("Email already exists")
	}

	phoneUser, err := u.userRepo.FindUserByPhone(req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if phoneUser != nil {
		return nil, errors.New("Phone already exists")
	}

	user := &entities.User{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	u.userRepo.CreateUser(user)

	return &RegisterResponse{
		ID:       user.Id,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}
