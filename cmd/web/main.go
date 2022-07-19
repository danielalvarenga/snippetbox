package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// * Create a application struct to inject dependencies needed for the whole application
// * For handlers spread in different packages we can pass dependencies directly to handler
// functions using closure (ex: https://www.calhoun.io/5-useful-ways-to-use-closures-in-go)
// or creating a handler struct for each package to receive dependencies to be accessed by
// the handler functions as methods
type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	// * Map a flag passed from commands in terminal to a variable, storing the default value,
	// such as `go run ./cmd/web -addr=":9999"`
	// * Return a pointer
	// * flag.String(), flagInt(), flag.Bool(), flag.Float64() try to convert value automatically
	// * The default for boolean flags is true
	// * When use the flag "-help" on the terminal all the mapped flags will be showed
	// * To map flags in pre-existing variables use `flag.StringVar(&addr, "addr", ":4000", "...")`
	// * Instead flags you can use env variables to get values with `os.Getenv()`, but for this you
	// don't have implicity type conversion, default value and help like flags. Another way is to use
	// flags with envs to have all the functionalities. Ex: `go run ./cmd/web -addr=$MY_ADDR_ENV`
	addr := flag.String("addr", ":4000", "HTTP network address")
	// * Parse flags to set the received values in mapped variables
	flag.Parse()

	// * With `log.New(destination, prefix, additional information)` we can set different log levels
	// for diferent destinations and different additional information (flags)
	// * Use log constants (https://pkg.go.dev/log#pkg-constants) to set additional information like
	// `log.Ldate` to add date, `log.Ltime` to add time, `log.Lshortfile` to add filename and line,
	// log.LUTC to use UTC datetime instead local datetime, etc...
	// * Loggers created by log.New() are concurrency-safe, but if you have many loggers writing in the
	// same destination you need to ensure that the destination is safe for concurrent use too
	// * Output: INFO    2022/07/19 19:32:12 my custom message
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// * Output: ERROR   2022/07/19 19:32:12 main.go:85: error message
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Instantiate the application injecting dependencies
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

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

	// The function http.ListenAndServe() should receive the TCP network
	// in the format "host:port" to specifics host and port number or
	// ":port" to any host and specific port number or
	// ":my-port-name" to use named port that Go will try to get the
	// correspondence from "/etc/services"
	infoLog.Println("Starting server on", *addr)
	// * The servemux implements the http.Handler interface, so we can pass it to the
	// http.ListenAndServe and for each request the handler method in servemux forward
	// to the correct handler method based in the registered routes
	// * Each request is handled concurrently in its own goroutine
	// * To use the custom errorLog for http errors we should to instantiate a new http
	// server struct passing the settings by properties and start the server from it
	srv := &http.Server{
		Addr:     *addr,
		Handler:  mux,
		ErrorLog: errorLog,
	}
	err := srv.ListenAndServe()
	// * Is a good practice call Panic() ou Fatal() and their variations only from the main()
	// * More log methods good to know: https://pkg.go.dev/log#Logger
	errorLog.Fatal(err)
}
