package application

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/handlers"
	m "github.com/punpundada/shelfMaster/internals/handlers/middleware"
	"github.com/punpundada/shelfMaster/internals/utils"
)

func loadRoutes(q *db.Queries) *chi.Mux {
	router := chi.NewRouter()
	mw := &m.Middleware{
		Queries: q,
	}
	// router.Use(mw.CSRFProtection)
	router.Use(mw.SetContentType)
	router.Use(mw.ValidateSessionCookie)
	router.Use(middleware.Logger)
	// router.Use(mw.TimeoutRequest)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserFromContext(r.Context())
		if err == nil {
			data, err := json.Marshal(&user)
			if err == nil {
				w.Write(data)
			}
		}
	})
	router.Route("/auth", loadAuthRoutes(q))
	router.Route("/admin", loadAdminRoutes(q))
	return router
}

func loadAuthRoutes(q *db.Queries) func(chi.Router) {
	authRoutess := handlers.NewAuth(q)

	return func(router chi.Router) {
		router.Post("/login", authRoutess.LoginUser)
		router.Post("/signup", authRoutess.RegisterUser)
	}
}

func loadAdminRoutes(q *db.Queries) func(chi.Router) {
	adminRoutes := handlers.NewAdmin(q)
	return func(router chi.Router) {
		router.Patch("/create/{id}", adminRoutes.CreateNewAdmin)
	}
}
