package server

import (
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTP struct {
	server *Router
}

type Router struct {
	*http.Server
	mux http.Handler
}

func createHandler(cfg Config) *Router {
	var mux *httprouter.Router = httprouter.New()
	addRoutes(mux)
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

func NewHttpServer(cfg Config) *HTTP {
	var srv *HTTP = &HTTP{
		server: createHandler(cfg),
	}
	return srv
}

func (srv *HTTP) Start() error {
	err := srv.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("error starting server: %s\n", err)
		return err
	}
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
