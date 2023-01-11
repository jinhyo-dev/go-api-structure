package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Get("/", app.Home)
	mux.Post("/auth", app.authentication)
	mux.Post("/sign-up", app.Register)
	mux.Delete("/delete-user", app.DeleteUser)
	mux.Get("/movies", app.AllMovies)
	return mux
}
