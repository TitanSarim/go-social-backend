package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("conflicting record")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Patch(context.Context, int64, *Post) error
	}
	Users interface {
		Create(context.Context, *User) (*User, error)
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
        GETByPostID(context.Context, int64) ([]Comment, error)
    }
    // Add other necessary interfaces for interacting with the database here
	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
    
	return Storage{
        Posts: &PostStore{db},
		Users: &UserStore{db},
		Comments: &CommentStore{db},
		Followers: &FollowerStore{db},
    }
}