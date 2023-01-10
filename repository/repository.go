package repository

import (
	"board-backend/models"
	"gorm.io/gorm"
)

type DatabaseRepo interface {
	Connection() *gorm.DB
	AllMovies() ([]*models.Movie, error)
}
