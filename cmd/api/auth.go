package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/TitanSarim/go-social-backend/internal/store"
	"github.com/google/uuid"
)

type CreateUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserWithToken struct{
	*store.User
	Token string `json:"token"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Email:    payload.Email,
		Username: payload.Username,
	}
	
	// hash user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
        return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp)
	if err != nil { 
		switch err{
			case store.ErrConflict:
				app.badRequestResponse(w, r, err)
		    default:
				app.internalServerError(w, r, err)
		}
	}

	userWithToken := UserWithToken{
		User:  user,
        Token: plainToken,
	}

	app.jsonResponse(w, http.StatusCreated, userWithToken)
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {}