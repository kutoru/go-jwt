package endpoints

import (
	"github.com/gorilla/mux"
)

// Creates a new router and sets up the routes
func LoadRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/auth", auth).Methods("POST")
	r.HandleFunc("/refresh", refresh).Methods("GET")

	return r
}
