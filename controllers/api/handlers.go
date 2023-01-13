package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
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

func (app *application) testReadJSON(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(requestPayload)
}

func (app *application) authentication(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	fmt.Println("requestPayload", requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	fmt.Println(requestPayload.Email)
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("user not found"), http.StatusBadRequest)
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

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	var payload JSONResponse
	payload.Success = true
	payload.Message = "auth success"
	payload.Data = tokens

	app.writeJSON(w, http.StatusAccepted, payload)
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

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
			}

			user, err := app.DB.GetUserById(userId)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating token"), http.StatusUnauthorized)
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))
			app.writeJSON(w, http.StatusOK, tokenPairs)
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}
