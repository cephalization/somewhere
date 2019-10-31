package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

// Global config object with parsed and validated program
var config map[string]interface{}

// Print usage message and then exit
func usage() {
	fmt.Printf("Usage: somewhere [arguments] directory\n\n")
	flag.PrintDefaults()
	os.Exit(1)
}

// Parse and map arguments for later usage
func setupArgs() (map[string]interface{}, error) {
	proxyPrefix := flag.String("prefix", "api", "route prefix that will be proxied. All other routes will be served the SPA")
	proxyScheme := flag.String("pscheme", "http://", "target host scheme to proxy api requests to (ex. 'https://')")
	proxyHost := flag.String("phost", "0.0.0.0", "target host to proxy api requests to")
	proxyPort := flag.String("pport", "8081", "target port to proxy api requests to")
	port := flag.String("port", "8080", "port to run server on")
	host := flag.String("host", "0.0.0.0", "host to run server on")

	flag.Parse()

	dir := flag.Arg(0)

	if dir == "" {
		usage()
	}

	mapping := map[string]interface{}{
		"port":         *port,
		"host":         *host,
		"directory":    dir,
		"proxy_host":   *proxyHost,
		"proxy_port":   *proxyPort,
		"proxy_scheme": *proxyScheme,
		"proxy_prefix": *proxyPrefix,
	}

	return mapping, nil
}

// Reverse proxy request to remote server
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	urlString :=
		config["proxy_scheme"].(string) +
			config["proxy_host"].(string) + ":" +
			config["proxy_port"].(string)
	u, err := url.Parse(urlString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)
}

// Serve single page application static files
func spaHandler(w http.ResponseWriter, r *http.Request) {
	// Load absolute path of the static directory
	dir, err := filepath.Abs(config["directory"].(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Combine static directory with the requested file
	path := filepath.Join(dir, r.URL.Path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// Return index.html, must be a client side route
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
		return
	} else if err != nil {
		// Some internal error has occurred
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serve the requested file, it exists
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

func main() {
	args, e := setupArgs()

	if e != nil {
		panic(e)
	}

	config = args

	r := mux.NewRouter()

	prefix := config["proxy_prefix"].(string)

	// Handle routes
	r.HandleFunc("/"+prefix+"/", proxyHandler)
	r.HandleFunc("/"+prefix+"/{route}", proxyHandler)
	r.PathPrefix("/").HandlerFunc(spaHandler)

	// Load config and start server
	host := config["host"].(string)
	port := config["port"].(string)
	srv := http.Server{
		Handler:      r,
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server listening at %s:%s", host, port)
	log.Fatal(srv.ListenAndServe())
}
