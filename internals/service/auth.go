package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type AuthService struct {
	Queries *db.Queries
}

type LoginBody struct {
	Email string `json:"email"`
}

func (a *AuthService) LoginUser(r *http.Request) (*db.Librarian, *db.Session, error) {
	var body = LoginBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, nil, err
	}

	if (len(body.Email) == 0) || (!utils.IsValidEmail(body.Email)) {
		return nil, nil, fmt.Errorf("invalid email")
	}
	librarian, err := a.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found: %v", err)
	}
	session, err := a.Queries.SaveSession(r.Context(), db.SaveSessionParams{
		UserID:    librarian.UserID,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 14)},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while saving session")
	}
	return &librarian, &session, nil
}

func (a *AuthService) SaveUser(ctx context.Context, user db.SaveUserParams) (db.User, error) {
	return a.Queries.SaveUser(ctx, user)
}
