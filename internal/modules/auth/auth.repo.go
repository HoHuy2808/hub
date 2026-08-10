package auth

import (
	"hub/internal/database/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	FindUserByPhone(phone string) (*entities.User, error)
	FindUserName(userName string) (*entities.User, error)
	CheckEmptyDB() bool
	CreateUser(user *entities.User) error
	CreateRefreshToken(refreshToken *entities.RefreshToken) error
}

type UserRepositoryImp struct {
	db *gorm.DB
}

func NewUserRepositoryImp(db *gorm.DB) UserRepository {
	return &UserRepositoryImp{db: db}
}

func (u *UserRepositoryImp) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := u.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) FindUserByPhone(phone string) (*entities.User, error) {
	var user entities.User
	result := u.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) FindUserName(userName string) (*entities.User, error) {
	var user entities.User
	result := u.db.Where("user_name = ?", userName).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImp) CheckEmptyDB() bool {
	var count int64
	u.db.Model(&entities.User{}).Count(&count)
	if count == 0 {
		return true
	}
	return false
}

func (u *UserRepositoryImp) CreateUser(user *entities.User) (err error) {
	result := u.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return err
}

func (u *UserRepositoryImp) CreateRefreshToken(refreshToken *entities.RefreshToken) error {
	result := u.db.Create(refreshToken)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
