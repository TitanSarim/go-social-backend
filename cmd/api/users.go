package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/TitanSarim/go-social-backend/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type FollowUser struct{
	UserID int64 `json:"user_id"`
}

type userKey string
const userCtx userKey = "user"

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
		Email: payload.Email,
        Username: payload.Username,
        Password: payload.Password,
	}

    ctx := r.Context()

    user, err := app.store.Users.Create(ctx, user)
    if err != nil {
        app.internalServerError(w, r, err)
        return
    }

    app.jsonResponse(w, http.StatusCreated, user)
    return
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request)    {


	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
        return
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {}

// user following
func (app *application) followHandler(w http.ResponseWriter, r *http.Request) {
	followUser := getUserFromCtx(r)
	if followUser == nil {
		app.badRequestResponse(w, r, errors.New("unauthorized or user not found"))
		return
	}

	//todo: revert it back to its original state
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

	ctx := r.Context()

	// this is line 85
	if err := app.store.Followers.Follow(ctx, int64(followUser.ID), payload.UserID); err != nil {
		app.internalServerError(w, r, err)
        return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
        return
	}
}

func (app *application) unFollowHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := getUserFromCtx(r)

	//todo: revert it back to its original state
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

	ctx := r.Context()

	if err := app.store.Followers.UnFollow(ctx, int64(unFollowedUser.ID), payload.UserID); err != nil {
		app.internalServerError(w, r, err)
        return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
        return
	}
}


func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		

		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func getUserFromCtx (r *http.Request) *store.User{
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}