package application

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/handlers"
	m "github.com/punpundada/shelfMaster/internals/handlers/middleware"
	"github.com/punpundada/shelfMaster/internals/utils"
)

func loadRoutes(q *db.Queries, Conn *pgx.Conn) *chi.Mux {
	router := chi.NewRouter()
	mw := &m.Middleware{
		Queries: q,
	}
	// router.Use(mw.CSRFProtection)
	router.Use(mw.SetContentType)
	router.Use(mw.ValidateSessionCookie)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// router.Use(mw.TimeoutRequest)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserFromContext(r.Context())
		if err != nil {
			w.Write([]byte(`{"no_user_found":"no_user_found"}`))
			return
		}
		data, _ := json.Marshal(user)
		w.Write(data)

	})
	router.Route("/auth", loadAuthRoutes(q, Conn))
	router.Route("/admin", loadAdminRoutes(q, mw))
	return router
}

func loadAuthRoutes(q *db.Queries, conn *pgx.Conn) func(chi.Router) {
	authRoutess := handlers.NewAuth(q, conn)

	return func(router chi.Router) {
		router.Post("/login", authRoutess.LoginUser)
		router.Post("/signup", authRoutess.RegisterUser)
		router.Post("/email-verification", authRoutess.EmailVerification)
		router.Post("/reset-password", authRoutess.ResetPassword)
		router.Post("/reset-password/{tokenId}", authRoutess.VeryfyRestPassword)
	}
}

func loadAdminRoutes(q *db.Queries, mw *m.Middleware) func(chi.Router) {
	adminRoutes := handlers.NewAdmin(q)
	return func(router chi.Router) {
		router.Use(mw.AdminOnly)
		router.Patch("/create/{id}", adminRoutes.CreateNewAdmin)
		router.Patch("/create/lib/{id}", adminRoutes.CreateLibrarian)
	}
}
