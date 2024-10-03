package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/service"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type Auth struct {
	service.AuthService
	Conntection *pgx.Conn
}

func NewAuth(q *db.Queries, conn *pgx.Conn) *Auth {
	return &Auth{
		AuthService: service.AuthService{
			Queries: q,
		},
		Conntection: conn,
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

func (a *Auth) EmailVerification(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Code string `json:"code"`
	}{}
	err := utils.ParseJSON(r, &body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "error decoding json"+err.Error())
		return
	}
	ses, err := utils.GetSessionFromContext(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid session: "+err.Error())
		return
	}
	_, user, err := utils.ValidateSession(r.Context(), a.Queries, ses.ID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "session not found: "+err.Error())
		return
	}
	if user == nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "user not found: ")
		return
	}
	tx, err := a.Conntection.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "error starting transtaction"+err.Error())
		return
	}
	isValidCode, err := utils.VerifyVerificationCode(r.Context(), tx, a.Queries, user, body.Code)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if !isValidCode {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "invalid code")
		return
	}

}
