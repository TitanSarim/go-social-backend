package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	IsActive bool `json:"is_active"`
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
	err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).
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

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// update the user's status to active
		user.IsActive = true
		if err := s.updateUserStatus(ctx, tx, user); err != nil {
			return err
		}

		// clean the invitation
		if err := s.deleteUserInvitation(ctx, tx, int64(user.ID)); err != nil {
			return err
		}

		return nil
	})
}


func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string,) (*User, error){
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users u
		JOIN user_invitations i ON u.id = i.user_id
		WHERE i.token = $1 AND i.expires_at > $2
		`
	hash := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx,  time.Second * 5)

	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)

	if err == sql.ErrNoRows {
        switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
		
    }

	return user, nil

}

func (s *UserStore) updateUserStatus(ctx context.Context, tx *sql.Tx, user *User) error {
	query := "UPDATE users SET is_active = true WHERE id = $1"

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.ID) 
	if err != nil {
		return err
	}
	return nil
}


func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error{
	query := "DELETE FROM user_invitations WHERE user_id = $1"

    ctx, cancel := context.WithTimeout(ctx,  time.Second * 5)

    defer cancel()

    _, err := tx.ExecContext(ctx, query, userID)
    if err != nil {
        return err
    }
    return nil    
}