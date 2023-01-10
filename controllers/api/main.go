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
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	var app application
	flag.StringVar(&app.DSN, "dsn", "root:1234@tcp(127.0.0.1:3306)/go_api?charset=utf8mb4&parseTime=True&loc=Local", "MariaDB Connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "singing secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "singing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "singing audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.Parse()

	conn, err := app.connectDatabase(app.DSN)
	app.DB = &dbrepo.MariaDBRepo{DB: conn}
	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
