package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Anant-raj2/tutorme/internal/auth"
	"github.com/Anant-raj2/tutorme/internal/tutor"
	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
	server       *Router
	tutorHandler *tutor.Handler
	authHandler  *auth.Handler
}

type Router struct {
	*http.Server
	mux *httprouter.Router
}

func createHandler(cfg Config) *Router {
	var mux *httprouter.Router = httprouter.New()
	var server *http.Server = &http.Server{
		Handler:      mux,
		Addr:         net.JoinHostPort(cfg.Host, cfg.Port),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	return &Router{
		server,
		mux,
	}
}

func NewHttpServer(cfg Config, tutorHandler *tutor.Handler, authHandler *auth.Handler) *HTTP {
	var srv *HTTP = &HTTP{
		server:       createHandler(cfg),
		tutorHandler: tutorHandler,
		authHandler:  authHandler,
	}
	srv.addRoutes(srv.server.mux)
	return srv
}

func (srv *HTTP) Start(ctx context.Context) error {
	go func() {
		log.Printf("listening on %s\n", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (thanks Alessandro Rosetti)
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := srv.server.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
