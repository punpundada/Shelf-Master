package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
)

type SessionService struct {
	Queries *db.Queries
}

func (a *AdminService) CreateSession(r *http.Request) (*db.Session, error) {
	var body db.SaveSessionParams
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error decoding body: %v", err)
	}
	session, err := a.Queries.SaveSession(r.Context(), body)
	if err != nil {
		return nil, fmt.Errorf("error saving user session: %v", err)
	}
	return &session, nil
}
