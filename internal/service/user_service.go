package service

import (
	"github.com/randnull/banner-service/pkg/models"
)

type UserRepostory interface {
	AddUser(register_form *models.Register, is_admin bool) error
	GetUser(username string, password string) (*models.User, error)
}

type UserService struct {
	repo UserRepostory
}

func NewUserSevice(repo UserRepostory) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) CreateUser(register_form models.Register, is_admin bool) error {
	err := service.repo.AddUser(&register_form, is_admin)

	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) GetUser(username string, password string) (*models.User, error) {
	user, err := service.repo.GetUser(username, password)

	if err != nil {
		return nil, err
	}

	return user, nil
}
