package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cephalization/somewhere/configutil"
	"github.com/cephalization/somewhere/service"
)

// Print usage message and then exit
func usage() {
	fmt.Printf("Usage: somewhere [arguments] directory\n\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	config, e := configutil.ParseConfig()

	if e != nil {
		usage()
	}

	if !config.Initialized {
		panic(errors.New("config not initialized. Exiting"))
	}

	server := service.Server{Config: config}

	server.Serve()
}
