package handlers

import (
	"net/http"

	"github.com/punpundada/libM/internals/service"
)

type Auth struct {
	service.AuthService
}

func (a *Auth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	a.SaveUser(r.Context(), email)
}
