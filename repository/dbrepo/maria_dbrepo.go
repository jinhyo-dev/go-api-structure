package dbrepo

import (
	"go-api-structure/models"
	"gorm.io/gorm"
	"time"
)

type MariaDBRepo struct {
	DB *gorm.DB
}

type Movies struct {
	ID          int
	Title       string
	ReleaseDate time.Time
	RunTime     int
	MPAARating  string
	Description string
	Image       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *MariaDBRepo) Connection() *gorm.DB {
	return m.DB
}

func (m *MariaDBRepo) AllMovies() ([]*models.Movie, error) {
	var movies []*models.Movie
	values := []Movies{
		{
			ID:          1,
			Title:       "Avatar",
			ReleaseDate: time.Now(),
			RunTime:     180,
			MPAARating:  "idk",
			Description: "Blue man",
			Image:       "PNG",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now()},
		{
			ID:          2,
			Title:       "HERO",
			ReleaseDate: time.Now(),
			RunTime:     180,
			MPAARating:  "idk",
			Description: "HERO",
			Image:       "JPG",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now()},
	}
	m.DB.AutoMigrate(&movies)
	m.DB.Create(&values)
	m.DB.Find(&movies).Table("movies").Order("title asc")
	return movies, nil
}
