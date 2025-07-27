package handlers

import (
	"encoding/json"
	"net/http"

	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/service"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type BooksHandler struct {
	service service.BookService
}

func NewBookservice(q *db.Queries) BooksHandler {
	handler := BooksHandler{
		service: service.NewBookService(q),
	}
	return handler
}

func (b *BooksHandler) SaveNewBook(w http.ResponseWriter, r *http.Request) {
	body := db.Book{}
	apiErr := utils.ParseJSON(r, &body)
	if apiErr != nil {
		apiErr.WriteError(w)
	}

	book, apiErr := b.service.SaveNewBook(r.Context(), &body)

	if apiErr != nil {
		apiErr.WriteError(w)
	}
	err := json.NewEncoder(w).Encode(&book)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
