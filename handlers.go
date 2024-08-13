package http

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/bnkamalesh/errors"
	"github.com/bnkamalesh/webgo/v7"

	"github.com/bnkamalesh/goapp/internal/api"
	"github.com/bnkamalesh/goapp/internal/pkg/logger"
)

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	apis api.Server
	home *template.Template
}

func (h *Handlers) routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:          "helloworld",
			Pattern:       "",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.HelloWorld)},
			TrailingSlash: true,
		},
		{
			Name:          "health",
			Pattern:       "/-/health",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.Health)},
			TrailingSlash: true,
		},
		{
			Name:          "create-user",
			Pattern:       "/users",
			Method:        http.MethodPost,
			Handlers:      []http.HandlerFunc{errWrapper(h.CreateUser)},
			TrailingSlash: true,
		},
		{
			Name:          "read-user-byemail",
			Pattern:       "/users/:email",
			Method:        http.MethodGet,
			Handlers:      []http.HandlerFunc{errWrapper(h.ReadUserByEmail)},
			TrailingSlash: true,
		},
	}
}

// Health is the HTTP handler to return the status of the app including the version, and other details
// This handler uses webgo to respond to the http request
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) error {
	out, err := h.apis.ServerHealth()
	if err != nil {
		return err
	}
	webgo.R200(w, out)
	return nil
}


