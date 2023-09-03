package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kutoru/go-jwt/db"
	"github.com/kutoru/go-jwt/endpoints"
)

func main() {
	// Loading the .env
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Could not load .env: %v\n", err)
	}

	// Connecting to the DB
	err = db.Connect(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Panicf("Could not connect to the DB: %v\n", err)
	}
	defer db.Close()

	// Resetting the DB
	// err = db.Reset()
	if err != nil {
		log.Panicf("Could not reset the DB: %v\n", err)
	}

	// Setting up the router
	r := endpoints.LoadRouter()
	http.Handle("/", &RouterDec{Router: r})

	// Starting the server
	port := ":4000"
	log.Printf("Listening on port %s\n", port)
	log.Panicln(http.ListenAndServe(port, nil))
}

type RouterDec struct {
	Router *mux.Router
}

func (routerDec *RouterDec) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(
		"Access-Control-Allow-Origin", r.Header.Get("Origin"),
	)
	w.Header().Set(
		"Access-Control-Allow-Credentials", "true",
	)
	w.Header().Set(
		"Access-Control-Allow-Methods",
		"GET, POST",
	)

	if r.Method == "OPTIONS" {
		return
	}

	routerDec.Router.ServeHTTP(w, r)
}
