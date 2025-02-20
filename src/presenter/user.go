package presenter

import (
	"w3st/models"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID    uuid.UUID   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserPresenter interface {
	ResponseUser(user *models.User) *UserResponse
}

type userPresenter struct{}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (u *userPresenter) ResponseUser(user *models.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}