package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (a *AuthService) LoginUser(r *http.Request) (*db.User, *db.Session, error) {
	var body = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&body)
	defer r.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	if (len(body.Email) == 0) || (!utils.IsValidEmail(body.Email)) {
		return nil, nil, fmt.Errorf("invalid email")
	}
	user, err := a.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("user not found either password or email do not match: %v", err)
	}
	isCorrectPassword := utils.VerifyHashString(user.PasswordHash, body.Password)
	if !isCorrectPassword {
		return nil, nil, fmt.Errorf("user not found either password or email do not match")
	}
	fmt.Println(user)
	sessionId := uuid.New()
	session, err := a.Queries.SaveSession(r.Context(), db.SaveSessionParams{
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 24 * 14)},
		ID:        sessionId.String(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while saving session: %v", err)
	}
	return &user, &session, nil
}

func (a *AuthService) SaveUser(ctx context.Context, user db.SaveUserParams) (db.User, error) {
	return a.Queries.SaveUser(ctx, user)
}
