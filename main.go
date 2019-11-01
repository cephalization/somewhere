package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cephalization/somewhere/configutil"
	"github.com/gorilla/mux"
)

var config *configutil.Config

// Print usage message and then exit
func usage() {
	fmt.Printf("Usage: somewhere [arguments] directory\n\n")
	flag.PrintDefaults()
	os.Exit(1)
}

// Reverse proxy request to remote server
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	urlString :=
		config.ProxyScheme +
			config.ProxyHost + ":" +
			config.ProxyPort
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
	dir, err := filepath.Abs(config.Directory)
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
	c, e := configutil.ParseConfig()

	if e != nil {
		usage()
	}
	if !c.Initialized {
		panic(errors.New("config not initialized. Exiting"))
	}

	config = c

	r := mux.NewRouter()

	prefix := config.ProxyPrefix

	// Handle routes
	r.HandleFunc("/"+prefix+"/", proxyHandler)
	r.HandleFunc("/"+prefix+"/{route}", proxyHandler)
	r.PathPrefix("/").HandlerFunc(spaHandler)

	// Load config and start server
	host := config.Host
	port := config.Port
	srv := http.Server{
		Handler:      r,
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server listening at %s:%s", host, port)
	log.Fatal(srv.ListenAndServe())
}
