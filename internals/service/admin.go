package service

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type AdminService struct {
	Queries *db.Queries
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := utils.GetUserFromContext(r.Context())
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		role, err := user.Role.Value()
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if role != db.RoleTypeADMIN {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (a *AdminService) CreateAdmin(ctx context.Context, id int32) (int32, error) {
	updatedId, err := a.Queries.CreateAdmin(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("error while updating user data: %v", err)
	}
	return updatedId, nil
}
