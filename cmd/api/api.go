package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ReyviRahman/social/internal/auth"
	"github.com/ReyviRahman/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type config struct {
	addr string
	db   dbConfig
	mail mailConfig
	auth authConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	exp time.Duration
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	authenticator auth.Authenticator
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	r.Route("/posts", func(r chi.Router) {
		r.Use(app.AuthTokenMiddleware)
		r.Post("/", app.createPostHandler)

		r.Route("/{postID}", func(r chi.Router) {
			r.Use(app.postsContextMiddleware)

			r.Get("/", app.getPostHandler)
			r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
			r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/activate/{token}", app.activateUserHandler)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/", app.getUserHandler)
			r.Put("/follow", app.followUserHandler)
			r.Put("/unfollow", app.unfollowUserHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Get("/feed", app.getUserFeedHandler)
		})
	})

	r.Route("/authentication", func(r chi.Router) {
		r.Post("/user", app.registerUserHandler)
		r.Post("/token", app.createTokenHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	app.logger.Infow("server has started", "addr", app.config.addr)

	return srv.ListenAndServe()
}
