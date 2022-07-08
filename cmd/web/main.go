package main

import (
	"log"
	"net/http"
)

func main() {
	// * Instead explicitly declare a servemux and use mux.HandleFunc(),
	// you can use http.HandleFunc() directly to set routes. It uses the global variable
	// "DefaultServeMux" implicitly created by "net/http" package, but it isn't safe
	// because any package can access global variables and handle this in your application.
	mux := http.NewServeMux()

	// * The Go file server serve static files
	// * Automatically sanitises paths (path.Clean()) to avoid directory traversal attacks
	// * Supports range requests to download large files
	// * If the file wasn't modified since the last request, it returns a "304 Not Modified"
	// status code instead the file
	// * If Content-Type is setted wrong (mime.TypeByExtension()) you can add your custom
	// type using the function mime.AddExtensionType()
	// * How to disable directory listing:
	// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
	// * To serve a single file direct from some handler function you can use http.ServeFile(), but
	// need to use filepath.Clean() to sanitize the path
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Create the route for static files removing the "/static" prefix
	// before the request reaches the file Server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Servemux supports two types of URL patterns:
	// * Subtree paths: it ends with trailing slash, like "/" or "/public/", and
	// catch any URl path starting with the pattern, like a wildcard ("/**" or "/public/**")
	// * Fixed paths: it doesn't end with trailing slash and is catched only when the
	// URL path is exactly the same
	//
	// Important:
	// * Servemux doesn't support RESTful functionalities like http methods, variables or
	// regexp in routes. Recommended to use third-party routers.
	// * Longer URL patterns takes precedence over shorter ones
	// * URL patterns are automatically sanitized and redirected (301 Permanent Redirect)
	// to the correspondent path. Ex: "/any/../foo//bar" to "/any/foo/bar"
	// * Request without trailing slash that matches with some subtree path is redirected
	// (301 Permanent Redirect) to subtree path. Ex: "/any" to "/any/"
	// * URL pattenrs accepts host names like "foo.any.com/create"
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// The function http.ListenAndServe() should receive the TCP network
	// in the format "host:port" to specifics host and port number or
	// ":port" to any host and specific port number or
	// ":my-port-name" to use named port that Go will try to get the
	// correspondence from "/etc/services"
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
