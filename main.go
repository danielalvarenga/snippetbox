package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// To avoid the servemux behavior for subtree paths,
	// check if the path is exactly the same to continue or redirect
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from SnippetBox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Servemux supports two types of URL patterns:
	// * Subtree paths: it ends with trailing slash, like "/" or "/public/", and
	// catch any URl path starting with the pattern, like a wildcard ("/**" or "/public/**")
	// * Fixed paths: it doesn't end with trailing slash and is catched only when the
	// URL path is exactly the same
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// The function http.ListenAndServe() should receive the TCP network
	// in the format "host:port" to specifics host and port number or
	// ":port" to any host and specific port number or
	// ":my-port-name" to use named port that Go will try to get the
	// correspondence from "/etc/services"
	log.Println("Stating server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
