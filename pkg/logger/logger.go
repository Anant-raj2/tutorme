package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func LogHttp(h func(w http.ResponseWriter, r *http.Request, _ httprouter.Params)) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var start time.Time = time.Now()
		h(w, r, ps)
		fmt.Printf("\nMETHOD: %s, TIME TAKEN: %d\n", r.Method, time.Since(start).Milliseconds())
	}
}
