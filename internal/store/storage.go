package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
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
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) (*User, error)
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
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

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error{
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
        return err
    }

	if err := fn(tx); err != nil{
		tx.Rollback()
		return err
	}

	return tx.Commit()
}