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

	"github.com/Anant-raj2/tutorme/internal/tutor"
	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
	server    *Router
	authStore *tutor.Store
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

func NewHttpServer(cfg Config, userHandler *tutor.Store) *HTTP {
	var srv *HTTP = &HTTP{
		server:    createHandler(cfg),
		authStore: userHandler,
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

// HelloWorld is a helloworld HTTP handler
func (h *Handlers) HelloWorld(w http.ResponseWriter, r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		{
			webgo.SendResponse(w, "hello world", http.StatusOK)
		}
	default:
		{
			buff := bytes.NewBufferString("")
			err := h.home.Execute(
				buff,
				struct {
					Message string
				}{
					Message: "Welcome to the Home Page!",
				},
			)
			if err != nil {
				return errors.InternalErr(err, "Inter server error")
			}

			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			_, err = w.Write(buff.Bytes())
			if err != nil {
				return errors.Wrap(err, "failed to respond")
			}
		}
	}
	return nil
}
