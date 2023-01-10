package repository

import (
	"go-api-structure/models"
	"gorm.io/gorm"
)

type DatabaseRepo interface {
	Connection() *gorm.DB
	AllMovies() ([]*models.Movie, error)
}
