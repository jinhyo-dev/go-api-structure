package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func (app *application) connectDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database connection is failed")
		return nil, err
	} else {
		log.Println("Connected to mysql!")
	}

	return db, nil
}
