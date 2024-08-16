package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func ErrorWrapper(h func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := h(w, r, nil)
		if err == nil {
			return
		}

		fmt.Fprintf(os.Stderr, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	})
}

func ParamsErrorWrapper(h func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		err := h(w, r, params)
		if err == nil {
			return
		}

		fmt.Fprintf(os.Stderr, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	})
}
