package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/punpundada/libM/internals/application"
	"github.com/punpundada/libM/internals/config"
	db "github.com/punpundada/libM/internals/db/sqlc"
)

func main() {
	ctx := context.Background()
	config := config.GetConfig()

	conn, err := pgx.Connect(ctx,
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			config.POSTGRES_USER,
			config.POSTGRES_PASSWORD,
			config.POSTGRES_HOST,
			config.POSTGRES_PORT,
			config.POSTGRES_DB))

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
