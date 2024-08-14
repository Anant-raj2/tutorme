package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ErrorWrapper(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	})
}
