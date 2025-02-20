package repositories

import (
	"w3st/models"

	"gorm.io/gorm"
)


type UserRepository interface {
	CreateUser(newUser *models.User) error
	FindUser(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}


// Create
func (r *userRepository) CreateUser(newUser *models.User) error {
	result := r.db.Create(newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Find
func (r *userRepository) FindUser (email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
	
}
