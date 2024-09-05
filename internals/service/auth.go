package service

import (
	db "github.com/punpundada/libM/internals/db/sqlc"
)

type AuthService struct {
	Queries *db.Queries
}

// func (a *AuthService) SaveUser(ctx context.Context, email string) (db.User, error) {

// }
