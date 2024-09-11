package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type Middleware struct {
	Queries *db.Queries
}

func (m *Middleware) CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			next.ServeHTTP(w, r)
		}
		originHeader := r.Header.Get("origin")
		hostHeader := r.Header.Get("host")
		if originHeader == "" || hostHeader == "" || utils.VerifyRequestOrigin(originHeader, []string{hostHeader}) {
			http.Error(w, "Forbidden - CSRF validation failed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type session string
type librarian string

const Sess session = "session"
const Libr librarian = "librarian"

func (m *Middleware) ValidateSessionCookie(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		session, librarian, err := utils.ValidateSession(r.Context(), m.Queries, cookie.Value)

		if err != nil {
			http.SetCookie(w, utils.CreateBlankSessionCookie())
			next.ServeHTTP(w, r)
			return
		}
		if session != nil && session.Fresh.Bool {
			_, err := m.Queries.UpdateSessionById(r.Context(), db.UpdateSessionByIdParams{
				ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 14), Valid: true},
				Fresh:     pgtype.Bool{Valid: true, Bool: true},
				ID:        session.ID,
			})
			if err != nil {
				http.SetCookie(w, utils.CreateBlankSessionCookie())
				next.ServeHTTP(w, r)
				return
			}
			http.SetCookie(w, utils.CreateSessionCookies(session.ID))
		}
		if session == nil {
			http.SetCookie(w, utils.CreateBlankSessionCookie())
		}
		contextWithData := context.WithValue(r.Context(), Sess, session)
		ctx := context.WithValue(contextWithData, Libr, librarian)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) SetContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
