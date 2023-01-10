package main

import (
	"flag"
	"fmt"
	"go-api-structure/repository"
	"go-api-structure/repository/dbrepo"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.DatabaseRepo
}

func main() {
	var app application
	flag.StringVar(&app.DSN, "dsn", "root:1234@tcp(127.0.0.1:3306)/go_api?charset=utf8mb4&parseTime=True&loc=Local", "MariaDB Connection")
	flag.Parse()
	conn, err := app.connectDatabase(app.DSN)
	app.Domain = "example.com"
	app.DB = &dbrepo.MariaDBRepo{DB: conn}
	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
