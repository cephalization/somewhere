package service

import (
	"errors"
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

// Server struct, instantiate with a config
type Server struct {
	Config *configutil.Config
}

func (s *Server) checkConfig() (*configutil.Config, error) {
	config := s.Config

	if !config.Initialized {
		return nil, errors.New("config not initialized, cannot serve")
	}

	return config, nil
}

// Reverse proxy request to remote server
func (s *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
	config, err := s.checkConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
func (s *Server) spaHandler(w http.ResponseWriter, r *http.Request) {
	config, err := s.checkConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

// Serve spa and reverse proxy api requests
func (s *Server) Serve() error {

	config := s.Config

	if !config.Initialized {
		return errors.New("config not initialized, cannot serve")
	}

	r := mux.NewRouter()

	prefix := config.ProxyPrefix

	// Handle routes
	r.HandleFunc("/"+prefix+"/", s.proxyHandler)
	r.HandleFunc("/"+prefix+"/{route}", s.proxyHandler)
	r.PathPrefix("/").HandlerFunc(s.spaHandler)

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

	return nil
}
