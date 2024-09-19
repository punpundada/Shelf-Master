package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/service"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type Auth struct {
	service.AuthService
}

func NewAuth(q *db.Queries) *Auth {
	return &Auth{
		AuthService: service.AuthService{
			Queries: q,
		},
	}
}

func (a *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(utils.SendVerificationEmail("prajwalparashkar100@gmail.com"))
	w.Write([]byte(`{"message":"This is Signup user route"}`))
}

func (a *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
	user, session, err := a.AuthService.LoginUser(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	http.SetCookie(w, utils.CreateSessionCookies(session.ID))
	if err = json.NewEncoder(w).Encode(struct {
		Email string `json:"email"`
	}{Email: user.Email}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func generateEmailVerificationCode(ctx context.Context, userId int32, email string, q *db.Queries) (string, error) {
	_, err := q.DeleteEmailVerificationByUserId(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("error deleting verifications: %v", err)
	}

	code := utils.GenerateRandomDigits(6)
	_, err = q.SaveEmailVerification(ctx, db.SaveEmailVerificationParams{
		Code:      code,
		UserID:    userId,
		Email:     email,
		ExpiresAt: pgtype.Date{Time: time.Now().Add(15 * time.Minute), Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("error saving varification code: %v", err)
	}
	return code, nil
}
