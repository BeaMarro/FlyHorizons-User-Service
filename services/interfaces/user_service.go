package interfaces

import (
	"flyhorizons-userservice/models"
)

type UserService interface {
	GetAll() []models.User
	GetByID(id int) (*models.User, error)
	UserExists(id int) bool
	Create(user models.User) (*models.User, error)
	DeleteByID(id int) (bool, error)
	Update(user models.User) (*models.User, error)
}
