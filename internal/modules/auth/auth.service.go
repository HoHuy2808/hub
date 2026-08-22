package auth

import (
	"context"
	"errors"
	"hub/internal/database/entities"
	"hub/pkg/authenticate"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

type UserServiceImp struct {
	userRepo UserRepository
}

func NewUserServiceImp(repo UserRepository) *UserServiceImp {
	return &UserServiceImp{userRepo: repo}
}

func (u *UserServiceImp) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	emailUser, err := u.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if emailUser != nil {
		return nil, errors.New("Email already exists")
	}

	phoneUser, err := u.userRepo.FindUserByPhone(ctx, req.Phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if phoneUser != nil {
		return nil, errors.New("Phone already exists")
	}

	usernameUser, err := u.userRepo.FindUserName(ctx, req.UserName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if usernameUser != nil {
		return nil, errors.New("Username already exists")
	}

	cost := bcrypt.DefaultCost
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		return nil, err
	}
	user := &entities.User{
		UserName: req.UserName,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(passwordHash),
	}
	if u.userRepo.CheckEmptyDB(ctx) == true {
		user.Role = "ADMIN"
	}
	err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return &RegisterResponse{
		ID:       user.Id,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

func (u *UserServiceImp) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Incorrect email or password")
	}

	isMatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if isMatch != nil {
		return nil, errors.New("Incorrect email or password")
	}

	AccessToken, err := authenticate.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	RefreshToken, err := authenticate.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := &entities.RefreshToken{
		RefreshToken: RefreshToken,
		UserId:       user.Id,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	err = u.userRepo.CreateRefreshToken(ctx, refreshTokenEntity)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		ID:           user.Id,
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
	}, nil
}
