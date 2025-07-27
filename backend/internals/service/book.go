package service

import (
	"context"
	"net/http"

	db "github.com/punpundada/shelfMaster/internals/db/sqlc"
	"github.com/punpundada/shelfMaster/internals/utils"
)

type BookService struct {
	queries *db.Queries
}

func NewBookService(queries *db.Queries) BookService {
	return BookService{
		queries: queries,
	}
}

func (b *BookService) SaveNewBook(ctx context.Context, book *db.Book) (*db.Book, *utils.ApiError) {
	savedBook, err := b.queries.SaveBookQuery(ctx, db.SaveBookQueryParams{
		Name:          book.Name,
		Authorid:      book.Authorid,
		Description:   book.Description,
		CreatedAt:     book.CreatedAt,
		UpdatedAt:     book.UpdatedAt,
		Isbn10:        book.Isbn10,
		Isbn13:        book.Isbn13,
		Edition:       book.Edition,
		Publisher:     book.Publisher,
		DatePublished: book.DatePublished,
	})
	if err != nil {
		return nil, utils.NewApiError(err.Error(), http.StatusInternalServerError)
	}
	return &savedBook, nil
}
