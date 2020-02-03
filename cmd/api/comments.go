package main

import (
	"net/http"

	"github.com/corymhall/blog-backend-go/cmd/api/models"
	ht "github.com/corymhall/blog-backend-go/pkg/http"
	"github.com/corymhall/blog-backend-go/pkg/render"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

func commentRouter(env *Env) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		env.CreateComment(w, rq, logger)
	})

	return r
}

func (env *Env) CreateComment(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger) {

	data := &CommentPayload{}
	if err := render.Bind(rq, data); err != nil {
		render.Render(w, rq, render.ErrInvalidRequest(err), logger)
		return
	}

	comment := data.Comment
	if err := env.db.DBCreateComment(comment); err != nil {
		render.Render(w, rq, render.ErrInternalServerError(err), logger)
		return
	}

	render.Status(rq, http.StatusCreated)
	render.Render(w, rq, NewCommentPayloadResponse(comment, env), logger)

}

type CommentPayload struct {
	Comment *models.Comment `json:"comment"`
	Replies []*ReplyPayload `json:"replies,omitempty"`
}

type CommentListPayload []*CommentPayload

func (p *CommentPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *CommentPayload) Bind(r *http.Request) error {
	return nil
}

func NewCommentPayloadResponse(comment *models.Comment, env *Env) *CommentPayload {
	resp := &CommentPayload{
		Comment: comment,
	}

	if resp.Replies == nil {
		if replies, _ := env.db.DBGetReplies(comment.PostID, comment.ID); replies != nil {
			resp.Replies = NewReplyListPayloadResponse(replies)
		}
	}

	return resp
}

func NewCommentListPayloadResponse(comments []*models.Comment, env *Env) []*CommentPayload {
	list := []*CommentPayload{}
	for _, comment := range comments {
		list = append(list, NewCommentPayloadResponse(comment, env))
	}
	return list
}
