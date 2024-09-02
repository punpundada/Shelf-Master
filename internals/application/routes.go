package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/punpundada/libM/internals/db/sqlc"
	"github.com/punpundada/libM/internals/handlers"
	"github.com/punpundada/libM/internals/service"
)

func loadRoutes(q *db.Queries) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world\n"))
	})
	router.Route("/auth", loadAuthRoutes(q))
	return router
}

func loadAuthRoutes(q *db.Queries) func(chi.Router) {
	authRoutess := handlers.Auth{
		AuthService: service.AuthService{
			Queries: q,
		},
	}

	return func(router chi.Router) {
		router.Post("/login", authRoutess.RegisterUser)
	}
}
