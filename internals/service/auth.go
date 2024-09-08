package service

import (
	"context"

	db "github.com/punpundada/libM/internals/db/sqlc"
)

type AuthService struct {
	Queries *db.Queries
}

func (a *AuthService) SaveUser(ctx context.Context, user db.SaveUserParams) (db.User, error) {
	return a.Queries.SaveUser(ctx, user)
}
