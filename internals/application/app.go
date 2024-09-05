package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/punpundada/libM/internals/config"
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
	PORT := config.GlobalConfig.PORT
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: a.route,
	}
	url := fmt.Sprintf("http://localhost:%s", PORT)
	fmt.Printf("\nServer Listning on %s\n", url)
	return server.ListenAndServe()
}
