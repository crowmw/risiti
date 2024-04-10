package service

import (
	"errors"

	"github.com/crowmw/risiti/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetAny() (bool, error)
}

type userService struct {
	db *gorm.DB
}

func DefaultUserService(db *gorm.DB) UserService {
	return &userService{
		db,
	}
}

func (s *userService) Create(user model.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &model.User{}, err
	}

	newUser := model.User{
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	result := s.db.Create(&newUser)
	if result.Error != nil {
		return &model.User{}, nil
	}

	return &newUser, nil
}

func (s *userService) GetByEmail(email string) (*model.User, error) {
	user := model.User{}
	result := s.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return &model.User{}, result.Error
	}

	return &user, nil
}

func (s *userService) GetAny() (bool, error) {
	user := model.User{}
	result := s.db.Take(&user)
	if result.Error != nil {
		return false, nil
	}
	return true, nil
}
