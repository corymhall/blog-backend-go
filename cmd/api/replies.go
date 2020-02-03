package main

import (
	"net/http"

	"github.com/corymhall/blog-backend-go/cmd/api/models"
	ht "github.com/corymhall/blog-backend-go/pkg/http"
	"github.com/corymhall/blog-backend-go/pkg/render"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

func replyRouter(env *Env) chi.Router {
	r := chi.NewRouter()

	r.Post("/", func(w http.ResponseWriter, rq *http.Request) {
		logger := ht.Logger(rq)
		env.CreateReply(w, rq, logger)
	})

	return r
}

func (env *Env) CreateReply(w http.ResponseWriter, rq *http.Request, logger zerolog.Logger) {
	data := &ReplyPayload{}
	if err := render.Bind(rq, data); err != nil {
		render.Render(w, rq, render.ErrInvalidRequest(err), logger)
		return
	}

	reply := data.Reply
	if err := env.db.DBCreateReply(reply); err != nil {
		render.Render(w, rq, render.ErrInternalServerError(err), logger)
		return
	}

	render.Status(rq, http.StatusCreated)
	render.Render(w, rq, NewReplyPayloadResponse(reply), logger)
}

type ReplyPayload struct {
	Reply *models.Reply `json:"reply"`
}

type ReplyListResponse []*ReplyPayload

func (p *ReplyPayload) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (p *ReplyPayload) Bind(r *http.Request) error {
	return nil
}

func NewReplyPayloadResponse(reply *models.Reply) *ReplyPayload {
	resp := &ReplyPayload{
		Reply: reply,
	}

	return resp
}

func NewReplyListPayloadResponse(replies []*models.Reply) []*ReplyPayload {
	list := []*ReplyPayload{}
	for _, reply := range replies {
		list = append(list, NewReplyPayloadResponse(reply))
	}
	return list
}
