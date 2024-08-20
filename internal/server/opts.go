package server

import (
	"time"
)

type Config struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func panicRecoverer(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		webgo.R500(w, errors.DefaultMessage)

		logger.Error(r.Context(), fmt.Sprintf("%+v", p))
		fmt.Println(string(debug.Stack()))
	}()

	next(w, r)
}

// Handlers struct has all the dependencies required for HTTP handlers
type Handlers struct {
	apis api.Server
	home *template.Template
}
