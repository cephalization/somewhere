package main

import (
	"flag"
	"fmt"
	"os"
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
	port := flag.String("port", "8080", "port to run server on")
	host := flag.String("host", "0.0.0.0", "host to run server on")

	flag.Parse()

	dir := flag.Arg(0)

	if dir == "" {
		usage()
	}

	mapping := map[string]interface{}{"port": *port, "host": *host, "directory": dir}

	return mapping, nil
}

func scopeTest() {
	fmt.Println(config)
}

func main() {
	args, e := setupArgs()

	if e != nil {
		panic(e)
	}

	config = args

	scopeTest()
}
