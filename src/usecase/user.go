package usecase

import (
	"w3st/interfaces/repositories"
	"w3st/interfaces/services"
	"w3st/models"
	"w3st/presenter"
	"w3st/utils"
)

type UserUsecase interface {
	CreateUser(name string, email string, password string) (string, error)
	FindUser(email string) (string, error)
}

type userUsecase struct {
	userRepo repositories.UserRepository
	authService  services.AuthService
	userPresenter presenter.UserPresenter
}

func NewUserUsecase(userRepo repositories.UserRepository, authService services.AuthService, userPresenter presenter.UserPresenter) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		authService: authService,
		userPresenter: userPresenter,
	}
}

func (u *userUsecase) CreateUser(name string, email string, password string) (string, error) {
	newUser := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err := u.userRepo.CreateUser(newUser)
	if err != nil {
		return "", err
	}
	stringID, err := utils.UuidToString(newUser.ID)
	if err != nil {
		return "", err
	}
	token, err := u.authService.GenerateToken(stringID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) FindUser(email string) (string, error) {
	user, err := u.userRepo.FindUser(email)
	if err != nil {
		return "", err
	}
	response := u.userPresenter.ResponseUser(user)
	token, err := u.authService.GenerateToken(response.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}