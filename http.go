package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

)

// Config holds all the configuration required to start the HTTP server
type Config struct {
	Host string
	Port uint16

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration

	TemplatesBasePath string
	EnableAccessLog   bool
}

type HTTP struct {
	listener string
	server   *webgo.Router
}

// Start starts the HTTP server
func (h *HTTP) Start() error {
	h.server.Start()
	return nil
}

func (h *HTTP) Shutdown(ctx context.Context) error {
	err := h.server.Shutdown()
	if err != nil {
		return errors.Wrap(err, "failed shutting down HTTP server")
	}

	return nil
}

