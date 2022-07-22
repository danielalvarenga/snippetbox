package main

import "net/http"

func (app *application) routes() *http.ServeMux {
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
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	// * mux.HandleFunc() is equivalent to "mux.Handle("/", http.HandlerFunc(home))", that
	// converts the normal function "home" in a accepted http.Handler interface
	mux.Handle("/snippet/create", http.HandlerFunc(app.snippetCreate))
	// * Passing dependencies from main() directly to the handler function using "closure"
	// using http.HandlerFunc() in the handler as a example
	mux.Handle("/closure-example", closureExample(app))

	return mux
}
