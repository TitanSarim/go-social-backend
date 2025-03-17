package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
        GETByPostID(context.Context, int64) ([]Comment, error)
    }
    // Add other necessary interfaces for interacting with the database here
}

func NewPostgresStorage(db *sql.DB) Storage {
    
	return Storage{
        Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
    }
}