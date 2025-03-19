package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Follow struct {
	UserID int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	CreatedAt string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}


func (s *FollowerStore) Follow(ctx context.Context, followerId, userId int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()
    _, err := s.db.ExecContext(ctx, query, userId, followerId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505"{
			return ErrConflict
		}
	}
    return err  // return error if any
}

func (s *FollowerStore) UnFollow(ctx context.Context, followerId, userId int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()
    _, err := s.db.ExecContext(ctx, query, userId, followerId)
    return err  // return error if any
}
