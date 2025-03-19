package main

import (
	"log"
	"net/http"
	"time"

	"github.com/TitanSarim/go-social-backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store store.Storage
}

type config struct {
	addr string
	db    dbConfig
	env   string
}

type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux{
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)
            r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
                r.Get("/", app.getUserHandler)
                r.Patch("/", app.updateUserHandler)
                r.Delete("/", app.deleteUserHandler)
				// user followers
				r.Put("/follow", app.followHandler)
				r.Put("/unfollow", app.unFollowHandler)
            })
		})
	})
	
	return r
}

func (app *application) run(mux *chi.Mux) error {

	srv := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("server has started at localhost%s", app.config.addr)

	return srv.ListenAndServe()
}