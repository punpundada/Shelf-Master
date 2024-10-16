package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/punpundada/shelfMaster/internals/config"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
)

type App struct {
	route http.Handler
}

func New(q *db.Queries, conn *pgx.Conn) *App {
	return &App{
		route: loadRoutes(q, conn),
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
