package httpx

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Err            error    `json:"-"` // low level error message
	HTTPStatusCode int      `json:"-"` // HTTP status code
	ErrorMessage   *Details `json:"error"`
}

type Details struct {
	StatusText  string `json:"status"`
	AppCode     int64  `json:"code,omitempty"`
	MessageText string `json:"message,omitempty"`
}

func (e *Response) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func BadRequest(err error) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		ErrorMessage: &Details{
			AppCode:     http.StatusBadRequest,
			StatusText:  http.StatusText(http.StatusBadRequest),
			MessageText: err.Error(),
		},
	}
}
