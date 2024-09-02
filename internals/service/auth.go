package service

import (
	"context"

	db "github.com/punpundada/libM/internals/db/sqlc"
)

type AuthService struct {
	Queries *db.Queries
}

func (a *AuthService) SaveUser(ctx context.Context, email string) (db.User, error) {
	return a.Queries.GetUserByEmail(ctx, email)
}
