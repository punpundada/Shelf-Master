package application

import (
	"context"
	"net/http"

	db "github.com/punpundada/libM/internals/db/sqlc"
)

type App struct {
	route http.Handler
}

func New(q *db.Queries) *App {
	return &App{
		route: loadRoutes(q),
	}
}

func (a *App) Start(cxt context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.route,
	}
	return server.ListenAndServe()
}
