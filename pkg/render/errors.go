package render

import (
	"net/http"
)

// ErrResponse Renderer type to handle all sorts of errors
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}


// Render method that sets the status code on the request.
// This is executed as part of the renderer function before the Response
// is written.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	Status(r, e.HTTPStatusCode)
	return nil
}

// returns a Renderer object that represents an invalid request error
func ErrInvalidRequest(err error) Renderer {
	return &ErrResponse{
		Err: err,
		HTTPStatusCode: 400,
		StatusText: "Invalid request.",
		ErrorText: err.Error(),
	}
}

// returns a Renderer object that represents an error rendering the response
func ErrRender(err error) Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

func ErrInternalServerError(err error) Renderer {
	return &ErrResponse{
		Err: err,
		HTTPStatusCode: 500,
		StatusText: "Internal server error.",
		ErrorText: err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

