package dbrepo

import (
	"errors"
	"fmt"
	"go-api-structure/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
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

func (m *MariaDBRepo) AddUser(userInformation models.UserSignUp) (bool, error) {

	m.DB.AutoMigrate(&models.User{})

	var exist bool
	m.DB.Raw("select true from users where email = ?", userInformation.Email).Scan(&exist)

	if exist {
		return false, nil
	} else {
		password, err := bcrypt.GenerateFromPassword([]byte(userInformation.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			return false, err
		}
		values := models.User{
			Email:     userInformation.Email,
			Password:  string(password),
			FirstName: userInformation.FirstName,
			LastName:  userInformation.LastName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		m.DB.Create(&values)
	}

	return true, nil
}

func (m *MariaDBRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	fmt.Println(email)
	m.DB.Raw("select * from users where email = ?", email).Scan(&user)
	if len(user.Email) == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (m *MariaDBRepo) DeleteUserById(userId int) (bool, error) {
	var exist bool
	m.DB.Raw("select true from users where id = ?", userId).Scan(&exist)

	if !exist {
		return false, errors.New("id is not exist")
	}

	m.DB.Raw("delete from users where id = ?", userId).Scan(&exist)
	return true, nil
}

func (m *MariaDBRepo) GetUserById(id int) (*models.User, error) {
	var user models.User
	m.DB.Find(&user).Where("id = ?", id)
	return &user, nil
}
