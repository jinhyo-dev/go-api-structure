package repository

import (
	"go-api-structure/models"
	"gorm.io/gorm"
)

type DatabaseRepo interface {
	Connection() *gorm.DB
	AllMovies() ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	AddUser(userInformation models.UserSignUp) (bool, error)
	DeleteUserById(userId int) (bool, error)
	GetUserById(id int) (*models.User, error)
}
