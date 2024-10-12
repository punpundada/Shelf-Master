package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/punpundada/shelfMaster/internals/config"
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
	if err := utils.ParseJSON(r, &body); err != nil {
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
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid code")
		return
	}
	_, err = a.Queries.UpdateUsersEmail_verification(r.Context(), db.UpdateUsersEmail_verificationParams{
		EmailVerified: pgtype.Bool{Bool: true, Valid: true},
		ID:            user.ID,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error updatin users: "+err.Error())
		return
	}
	session, err := a.Queries.SaveSession(r.Context(), *utils.NewSaveSessionAttrs(user.ID))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Error saving session: "+err.Error())
		return
	}
	sessionCookie := utils.CreateSessionCookies(session.ID)
	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, sessionCookie)
}

func (a *Auth) ResetPassword(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email string `json:"email"`
	}{}
	if err := utils.ParseJSON(r, &body); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Unable to parse body "+err.Error())
		return
	}
	user, err := a.Queries.GetUserByEmail(r.Context(), body.Email)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid email "+err.Error())
		return
	}
	verificationToken, err := utils.CreatePasswordRestToken(r.Context(), a.Queries, user.ID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	verificationLink := config.GetConfig().FRONTEND_URL + "reset-password/" + verificationToken
	err = utils.SendPasswordResetEmail(user.Email, verificationLink)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "unable to send email "+err.Error())
		return
	}
	fmt.Println("token", verificationToken)
	utils.WriteResponse(w, http.StatusOK, "Password resend link send", true)
}

func (a *Auth) VeryfyRestPassword(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Password string `json:"password"`
	}{}
	if err := utils.ParseJSON(r, &body); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Error parsing body "+err.Error())
		return
	}
	verificationToken := r.PathValue("tokenId")
	tokenHash := utils.EncodeString(verificationToken)
	fmt.Println("tokenHash", tokenHash)
	resetToken, err := a.Queries.GetResetPasswordFromTokenHash(r.Context(), pgtype.Text{
		String: tokenHash,
		Valid:  true,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request : "+err.Error())
		return
	}
	if !utils.IsWithinExpirationDate(resetToken.ExpiresAt.Time) {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "No response")
		return
	}
	if err = utils.InvalidateAllUserSessions(r.Context(), a.Queries, resetToken.UserID); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	passwordHash, err := utils.HashString(body.Password)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_, err = a.Queries.UpdateUserPasswordByUserId(r.Context(), passwordHash)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	session, err := a.Queries.SaveSession(r.Context(), *utils.NewSaveSessionAttrs(resetToken.UserID))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	sessionCookie := utils.CreateSessionCookies(session.ID)
	http.SetCookie(w, sessionCookie)
	w.Header().Add("Referrer-Policy", "strict-origin")
}
