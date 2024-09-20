package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	user, err := a.AuthService.SaveUser(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	session, err := a.Queries.SaveSession(r.Context(), db.SaveSessionParams{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(15 * time.Minute), Valid: true},
	})

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	sessionCooke := utils.CreateSessionCookies(session.ID)
	http.SetCookie(w, sessionCooke)
	err = json.NewEncoder(w).Encode(struct {
		UserId int32 `json:"user_id"`
	}{UserId: user.ID})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
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
