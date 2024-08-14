package handler

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HelloWorld(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	json.NewEncoder(w).Encode("Hello World")
	// return nil
}
