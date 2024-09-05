package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/punpundada/libM/internals/application"
	"github.com/punpundada/libM/internals/config"
	db "github.com/punpundada/libM/internals/db/sqlc"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, config.GlobalConfig.CONNECTION_STR)

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	queries := db.New(conn)

	app := application.New(queries)
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

}
