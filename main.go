package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display the home page")
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet: %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display a form for creating a new snippet")
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	w.WriteHeader(http.StatusCreated)

	fmt.Fprintf(w, "Post Handler for creating a new snippet")
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on http://localhost:4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
