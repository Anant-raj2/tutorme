package server

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Config struct {
	Host         string
	Port         uint16
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTTP struct {
	// Handler      http.Handler
	// Addr         string
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
	*http.Server
}

func NewServer(cfg *Config) *http.Server {
	var mux *httprouter.Router = httprouter.New()
	var listener string = cfg.Host
	return &http.Server{
		Handler:      mux,
		Addr:         listener,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

// func (srv *HTTP) Start() {
//   err:=srv.ListenAndServe()
//   if err != nil
// }
