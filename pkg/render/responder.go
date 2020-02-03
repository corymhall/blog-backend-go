package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"context"

	"github.com/rs/zerolog"
)

var (
	StatusCtxKey = &contextKey{"Status"}
)

// Status sets a HTTP response status code hint into request context at any point
// during the request life-cycle. Before the Responder sends its response header
// it will check the StatusCtxKey
func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusCtxKey, status))
}

func Respond(w http.ResponseWriter, r *http.Request, v interface{}, logger zerolog.Logger) {
	entry := logger.Info()
	if err, ok := v.(error); ok {

		// we set a default error status response code if one hasn't been set.
		status, ok := r.Context().Value(StatusCtxKey).(int)
		if !ok {
			w.WriteHeader(400)
		} 

		// we log the error
		entry.Msg(fmt.Sprintf("http error: %s (code=%d)", err, status))

		JSON(w, r, v, logger)
		return
	}
	JSON(w, r, v, logger)
}

func JSON(w http.ResponseWriter, r *http.Request, v interface{}, logger zerolog.Logger) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		Respond(w, r, ErrInternalServerError(err), logger)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	w.Write(buf.Bytes())
}


// NoContent returns a HTTP 204 "No Content" response.
func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
