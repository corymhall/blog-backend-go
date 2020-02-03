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

func postRouter(env *Env) chi.Router {
	r := chi.NewRouter()

	r.With(env.PostsCtx).Get("/", func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		GetPosts(w, rq, logger, env)
	})

	r.Post("/", func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		env.CreatePost(w, rq, logger)
	})

	r.Route("/{postID}", func(r chi.Router) {
		r.Use(env.PostCtx)
		r.Get("/", func(w http.ResponseWriter, rq *http.Request) {
			logger := ht.Logger(rq)
			GetPost(w, rq, logger, env)
		})
		r.Put("/", func(w http.ResponseWriter, rq *http.Request) {
			logger := ht.Logger(rq)
			UpdatePost(w, rq, logger, env)
		})
	})

	return r
}

func (env *Env) PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		var post *models.Post
		var err error

		if postID := chi.URLParam(rq, "postID"); postID != "" {
			post, err = env.db.DBGetPost(postID)
		} else {
			render.Render(w, rq, render.ErrNotFound, logger)
			return
		}
		if err != nil {
			render.Render(w, rq, render.ErrNotFound, logger)
			return
		}
		ctx := context.WithValue(rq.Context(), "post", post)
		next.ServeHTTP(w, rq.WithContext(ctx))
	})
}
func (env *Env) PostsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		var posts []*models.Post
		var err error

		posts, err = env.db.DBGetPosts()
		if err != nil {
			render.Render(w, rq, render.ErrNotFound, logger)
			return
		}
		ctx := context.WithValue(rq.Context(), "posts", posts)
		next.ServeHTTP(w, rq.WithContext(ctx))
	})
}

func (env *Env) CreatePost(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger) {
	data := &PostPayload{}

	if err := render.Bind(rq, data); err != nil {
		render.Render(w, rq, render.ErrInvalidRequest(err), logger)
		return
	}

	post := data.Post
	if err := env.db.DBCreatePost(post); err != nil {
		render.Render(w, rq, render.ErrInternalServerError(err), logger)
		return
	}

	render.Status(rq, http.StatusCreated)
	render.Render(w, rq, NewPostPayloadResponse(post, env), logger)
}

func UpdatePost(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger, env *Env) {
	post := rq.Context().Value("post").(*models.Post)

	data := &PostPayload{
		Post: post,
	}
	if err := render.Bind(rq, data); err != nil {
		render.Render(w, rq, render.ErrInvalidRequest(err), logger)
		return
	}

	post = data.Post
	if err := env.db.DBUpdatePost(post); err != nil {
		render.Render(w, rq, render.ErrInternalServerError(err), logger)
		return
	}
	render.Render(w, rq, NewPostPayloadResponse(post, env), logger)

}

func GetPost(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger, env *Env) {
	post := rq.Context().Value("post").(*models.Post)

	if err := render.Render(w, rq, NewPostPayloadResponse(post, env), logger); err != nil {
		render.Render(w, rq, render.ErrRender(err), logger)
		return
	}

}
func GetPosts(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger, env *Env) {
	posts := rq.Context().Value("posts").([]*models.Post)

	if err := render.RenderList(w, rq, NewPostListPayloadResponse(posts, env), logger); err != nil {
		render.Render(w, rq, render.ErrRender(err), logger)
		return
	}

}

type PostPayload struct {
	Post     *models.Post      `json:"post"`
	Comments []*CommentPayload `json:"comments,omitempty"`
}

type PostListResponse []*PostPayload

func (p *PostPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *PostPayload) Bind(r *http.Request) error {
	return nil
}

func NewPostPayloadResponse(post *models.Post, env *Env) *PostPayload {
	resp := &PostPayload{
		Post: post,
	}

	if resp.Comments == nil {
		if comments, _ := env.db.DBGetComments(post.ID); comments != nil {
			resp.Comments = NewCommentListPayloadResponse(comments, env)
		}
	}

	return resp
}

func NewPostListPayloadResponse(posts []*models.Post, env *Env) []render.Renderer {
	list := []render.Renderer{}
	for _, post := range posts {
		list = append(list, NewPostPayloadResponse(post, env))
	}
	return list
}
