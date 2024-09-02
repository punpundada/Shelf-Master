package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/punpundada/libM/internals/application"
	db "github.com/punpundada/libM/internals/db/sqlc"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://user:password@localhost:5432")
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
