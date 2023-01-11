package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json: "status"`
		Message string `json: "message"`
		Version string `json: "version"`
	}{
		Status:  "active",
		Message: "Go api is running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) authentication(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credential"), http.StatusBadRequest)
		return
	}

	u := jwtUser{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName}

	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	log.Println(tokens.Token)
	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	if len(requestPayload.Email) < 5 || len(requestPayload.Password) < 3 {
		app.errorJSON(w, errors.New("email and password must be at least 4 characters long"), http.StatusBadRequest)
	}

	success, err := app.DB.AddUser(requestPayload)

	if err != nil || !success {
		app.errorJSON(w, errors.New("email is already exist"), http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusAccepted, "register complete")
}

func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userid")
	id, _ := strconv.Atoi(userId)
	_, err := app.DB.DeleteUserById(id)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
	}

	app.writeJSON(w, http.StatusAccepted, "user is successfully deleted")
}
