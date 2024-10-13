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
			err.WriteError(w, "Forbidden")
			return
		}
		role, error := user.Role.Value()
		if error != nil {
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

func (a *AdminService) CreateLibrarian(ctx context.Context, userId int32) (int32, *utils.ApiError) {
	//update user with role = 'librarian'
	//runtuen id and error
	id, err := a.Queries.CreateLibrarian(ctx, userId)
	if err != nil {
		return 0, utils.NewApiError(err.Error(), http.StatusInternalServerError)
	}
	return id, nil
}
