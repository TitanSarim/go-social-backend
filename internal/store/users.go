package store

import (
	"context"
	"database/sql"
)

type User struct{
	ID    int    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}


func (s *UserStore) Create(ctx context.Context, user *User) (*User, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at"
	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := "SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1"
	
	user := &User{}
    err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
        return nil, err
    }
	return user, nil
   
}