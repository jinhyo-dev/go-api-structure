package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(app.enableCORS)
		r.Get("/", app.Home)
		r.Get("/refresh", app.refreshToken)
		r.Post("/auth", app.authentication)
		r.Post("/test", app.testReadJSON)
		r.Post("/sign-up", app.Register)
		r.Delete("/delete-user", app.DeleteUser)
		r.Get("/logout", app.logout)
		r.Get("/movies", app.AllMovies)

		r.Route("/admin", func(r chi.Router) {
			r.Use(app.authRequired)
			r.Get("/movies", app.MovieCatalog)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	return r
}
