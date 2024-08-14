package server

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func CreateService() (*Server, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	var srv *Server = &Server{
		Host:         os.Getenv("HOST"),
		Port:         os.Getenv("PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv, nil
}

func (srv *Server) createHttpServer() *http.Server {
	var mux *httprouter.Router = httprouter.New()
	addRoutes(mux)
	return &http.Server{
		Handler:      mux,
		Addr:         net.JoinHostPort(srv.Host, srv.Port),
		ReadTimeout:  srv.ReadTimeout,
		WriteTimeout: srv.WriteTimeout,
	}
}

func (srv *Server) Start() error {
	err := srv.createHttpServer().ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		return err
	}
	return nil
}
