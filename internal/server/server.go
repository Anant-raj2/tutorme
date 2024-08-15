package server

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	http *http.Server
}

func CreateService(cfg Config) *Server {
	var mux *httprouter.Router = httprouter.New()

	addRoutes(mux)

	var server *http.Server = &http.Server{
		Handler:      mux,
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	var srv *Server = &Server{
		http: server,
	}
	return srv
}

func (srv *Server) Start() error {
	err := srv.http.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		return err
	}
	return nil
}
