package handlers

import (
	"encoding/json"
	"net/http"

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
	w.Write([]byte(`{"message":"This is Signup user route"}`))
}

func (a *Auth) LoginUser(w http.ResponseWriter, r *http.Request) {
	lib, session, err := a.AuthService.LoginUser(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	http.SetCookie(w, utils.CreateSessionCookies(session.ID))
	if err = json.NewEncoder(w).Encode(struct {
		Email string `json:"email"`
	}{Email: lib.Email}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
