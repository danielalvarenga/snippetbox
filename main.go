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
	if r.Method != "POST" {
		// * You must set headers, if you need, before write response ("w.WriteHeader()",
		// "w.Write()" or http.Error())
		// * Use http constants is the best practice: http.MethodPost = "POST". More
		// contsants in https://pkg.go.dev/net/http#pkg-constants
		w.Header().Set("Allow", http.MethodPost)

		// * For errors use the helper "http.Error(...)" that calls "w.WriteHeader" and "w.Write"
		// indirectly to write the response;
		// * Use http constants is the best practice: http.StatusMethodNotAllowed = 405. More
		// contsants in https://pkg.go.dev/net/http#pkg-constants
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// * Go automatically set the headers "Content-type", "Content-Length" and "Date", but
	// the function http.DetectContentType() can't distinguish JSON from plain text, so
	// you need to set manually.
	w.Header().Set("Content-Type", "application/json")

	// To supress sytem generated headers you need to set nil for the value
	w.Header()["Date"] = nil

	// * Instead use "Set" you can use "Add" to append values for the same header
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-Control", "max-age=32645000")

	// * When editing the header throught Set(), Add(), Del(), Get() and Values(), it
	// will be automatically canonicalized (first letters in keys names putted in upper case for
	// HTTP/1 or all keys and values in down case for HTTP/2)
	w.Header()["down-case"] = []string{"true"}

	// * The function "w.WriteHeader" should be called only once per response;
	// * When not called after the function "w.Write", the response status code will be 200;
	w.WriteHeader(201)
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Servemux supports two types of URL patterns:
	// * Subtree paths: it ends with trailing slash, like "/" or "/public/", and
	// catch any URl path starting with the pattern, like a wildcard ("/**" or "/public/**")
	// * Fixed paths: it doesn't end with trailing slash and is catched only when the
	// URL path is exactly the same
	//
	// Important:
	// * Servemux doesn't support RESTful functionalities like http methods, variables or
	// regexp in routes. Recommended to use third-party routers (future commits).
	// * Longer URL patterns takes precedence over shorter ones
	// * URL patterns are automatically sanitized and redirected (301 Permanent Redirect)
	// to the correspondent path. Ex: "/any/../foo//bar" to "/any/foo/bar"
	// * Request without trailing slash that matches with some subtree path is redirected
	// (301 Permanent Redirect) to subtree path. Ex: "/any" to "/any/"
	// * URL pattenrs accepts host names like "foo.any.com/create"
	// * Instead explicitly declare a servemux and use mux.HandleFunc(),
	// you can use http.HandleFunc() directly using the global variable "DefaultServeMux"
	// implicitly created by "net/http" package, but it isn't safe because any package
	// can access global variables and handle this in your application
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
