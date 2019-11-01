package configutil

import (
	"errors"
	"flag"
)

// Config struct with parsed and validated program arguments
type Config struct {
	ProxyPrefix string
	ProxyScheme string
	ProxyHost   string
	ProxyPort   string
	Port        string
	Host        string
	Directory   string
	Initialized bool
}

// ParseConfig returns a new struct containing configuration parameters
func ParseConfig() (*Config, error) {
	proxyPrefix := flag.String("prefix", "api", "route prefix that will be proxied. All other routes will be served the SPA")
	proxyScheme := flag.String("pscheme", "http://", "target host scheme to proxy api requests to (ex. 'https://')")
	proxyHost := flag.String("phost", "0.0.0.0", "target host to proxy api requests to")
	proxyPort := flag.String("pport", "8081", "target port to proxy api requests to")
	port := flag.String("port", "8080", "port to run server on")
	host := flag.String("host", "0.0.0.0", "host to run server on")

	flag.Parse()

	dir := flag.Arg(0)

	if dir == "" {
		return nil, errors.New("Invalid usage")
	}

	c := Config{
		Port:        *port,
		Host:        *host,
		Directory:   dir,
		ProxyHost:   *proxyHost,
		ProxyPort:   *proxyPort,
		ProxyScheme: *proxyScheme,
		ProxyPrefix: *proxyPrefix,
		Initialized: true,
	}

	return &c, nil
}
