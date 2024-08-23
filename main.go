package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the home page"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet"))
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
