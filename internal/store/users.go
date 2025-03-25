package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct{
	ID    int    `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password password `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type password struct{
	text *string
	hash []byte
}

func (p *password) Set(text string) error{
	hash, err := bcrypt.GenerateFromPassword([] byte(text), bcrypt.DefaultCost)
	if err != nil {
        panic(err)
    }

	p.text = &text
	p.hash = hash
	return nil
}


type UserStore struct {
	db *sql.DB
}


func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) (*User, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at"
	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		switch{
			case err.Error() == `pq: duplicate key value violates unqiue constraint "users_email_key"`:
				return nil, ErrConflict
			case err.Error() == `pq: duplicate key value violates unqiue constraint "users_username_key"`:
				return nil, ErrConflict
			default:
				return nil, err
		}
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

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error{
		// create the user
		if _, err := s.Create(ctx, tx, user); err != nil{
			return err
		}

		// create user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, int64(user.ID)); err != nil{
			return err
		}

		return nil
	})
}

func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := "INSERT INTO user_invitations (token, user_id, expires_at) VALUES ($1, $2, $3)"
    _, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))
	if err != nil {
        return err
    }
    return nil
}