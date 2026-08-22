package auth

import (
	"context"
	"hub/internal/database/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entities.User, error)
	FindUserByPhone(ctx context.Context, phone string) (*entities.User, error)
	FindUserName(ctx context.Context, userName string) (*entities.User, error)
	CheckEmptyDB(ctx context.Context) bool
	CreateUser(ctx context.Context, user *entities.User) error
	CreateRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error
}

type UserRepositoryImp struct {
	db *gorm.DB
}

func NewUserRepositoryImp(db *gorm.DB) *UserRepositoryImp {
	return &UserRepositoryImp{db: db}
}

func (u *UserRepositoryImp) getDB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("TxKey").(*gorm.DB)
	if ok {
		return tx
	}
	return u.db
}

func (u *UserRepositoryImp) FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	db := u.getDB(ctx)
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) FindUserByPhone(ctx context.Context, phone string) (*entities.User, error) {
	var user entities.User
	db := u.getDB(ctx)
	result := db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) FindUserName(ctx context.Context, userName string) (*entities.User, error) {
	var user entities.User
	db := u.getDB(ctx)
	result := db.Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) CheckEmptyDB(ctx context.Context) bool {
	var count int64
	db := u.getDB(ctx)
	db.Model(&entities.User{}).Count(&count)
	if count == 0 {
		return true
	}
	return false
}

func (u *UserRepositoryImp) CreateUser(ctx context.Context, user *entities.User) (err error) {
	db := u.getDB(ctx)
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return err
}

func (u *UserRepositoryImp) CreateRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error {
	db := u.getDB(ctx)
	result := db.Create(refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
