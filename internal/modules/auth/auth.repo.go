package auth

import (
	"hub/internal/database/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
	FindUserByPhone(phone string) (*entities.User, error)
	FindUserName(userName string) (*entities.User, error)
	CreateUser(user *entities.User) error
}

type UserRepositoryImp struct {
	db *gorm.DB
}

func NewUserRepositoryImp(db *gorm.DB) *UserRepositoryImp {
	return &UserRepositoryImp{db: db}
}

func (r *UserRepositoryImp) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImp) FindUserByPhone(phone string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImp) FindUserName(userName string) (*entities.User, error) {
	var user entities.User
	result := r.db.Where("user_name = ?", userName).Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepositoryImp) CreateUser(user *entities.User) (err error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return err
}
