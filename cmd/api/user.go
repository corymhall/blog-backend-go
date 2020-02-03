package main

import (
	"context"
	"net/http"

	"github.com/corymhall/blog-backend-go/cmd/api/models"
	ht "github.com/corymhall/blog-backend-go/pkg/http"
	"github.com/corymhall/blog-backend-go/pkg/render"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

func userRouter(env *Env) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		env.CreateUser(w, rq, logger)
	})

	r.Route("/{userID}", func(r chi.Router) {
		r.Use(env.UserCtx)
		r.Get("/", func(w http.ResponseWriter, rq *http.Request) {
			logger := ht.Logger(rq)
			GetUser(w, rq, logger)
		})
	})

	return r
}

func (env *Env) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		var user *models.User
		var err error

		if userID := chi.URLParam(rq, "userID"); userID != "" {
			user, err = env.db.DBGetUser(userID)
		} else {
			render.Render(w, rq, render.ErrNotFound, logger)
			return
		}
		if err != nil {
			render.Render(w, rq, render.ErrNotFound, logger)
			return
		}
		ctx := context.WithValue(rq.Context(), "user", user)
		next.ServeHTTP(w, rq.WithContext(ctx))
	})
}

func (env *Env) CreateUser(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger) {

	data := &UserPayload{}
	if err := render.Bind(rq, data); err != nil {
		render.Render(w, rq, render.ErrInvalidRequest(err), logger)
		return
	}

	user := data.User
	if err := env.db.DBCreateUser(user); err != nil {
		render.Render(w, rq, render.ErrInternalServerError(err), logger)
		return
	}

	render.Status(rq, http.StatusCreated)
	render.Render(w, rq, NewUserPayloadResponse(user), logger)
}

func GetUser(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger) {
	user := rq.Context().Value("user").(*models.User)

	if err := render.Render(w, rq, NewUserPayloadResponse(user), logger); err != nil {
		render.Render(w, rq, render.ErrRender(err), logger)
		return
	}

}

type UserPayload struct {
	*models.User
}

func (p *UserPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *UserPayload) Bind(r *http.Request) error {
	return nil
}

func NewUserPayloadResponse(user *models.User) *UserPayload {
	resp := &UserPayload{user}

	return resp
}
