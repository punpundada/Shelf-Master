package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/service"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type Admin struct {
	Service service.AdminService
}

func NewAdmin(q *db.Queries) *Admin {
	return &Admin{
		Service: service.AdminService{
			Queries: q,
		},
	}
}

func (a *Admin) CreateNewAdmin(w http.ResponseWriter, r *http.Request) {

	pathId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error parsing id: %s", err.Error()))
		return
	}
	body := &struct {
		Id int32 `json:"id"`
	}{}
	err = json.NewDecoder(r.Body).Decode(body)
	defer r.Body.Close()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error parsing request body: %s", err.Error()))
		return
	}
	if pathId != int(body.Id) {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "error parsing request: Invalid id")
		return
	}
	updatedId, err := a.Service.CreateAdmin(r.Context(), int32(pathId))
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(utils.SuccessResponse{
		IsSuccess: true,
		Message:   "User promoted",
		Code:      http.StatusOK,
		Result:    updatedId,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error encoding: %s", err.Error()))
		return
	}
}

func (a *Admin) CreateLibrarian(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Id int32 `json:"id"`
	}{}
	if err := utils.ParseJSON(r, &body); err != nil {
		err.WriteError(w)
		return
	}
	fmt.Println("body", body)
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("invalid id: required integer passed string "+err.Error()))
		return
	}
	if id != int64(body.Id) {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "invalid id")
		return
	}
	returnedId, apiErr := a.Service.CreateLibrarian(r.Context(), int32(id))
	if apiErr != nil {
		apiErr.WriteError(w)
		return
	}
	body.Id = returnedId
	utils.WriteResponse(w, http.StatusOK, "User updated successfully", &body)

}
