package factory

import (
	"w3st/interfaces/controllers"
	"w3st/interfaces/repositories"
	"w3st/interfaces/services"
	"w3st/presenter"

	"w3st/usecase"

	"gorm.io/gorm"
)

type Factory struct {
	DB *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{DB: db}
}

func (f *Factory) NewUserController() *controllers.UserController {
	userRepo := repositories.NewUserRepository(f.DB)
	userPresenter := presenter.NewUserPresenter()
	authService := services.NewAuthService()
	userUsecase := usecase.NewUserUsecase(userRepo, authService, userPresenter)
	return controllers.NewUserController(userUsecase)
}
