package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// PORT - server port
var PORT = "8080"

// PATH - server path
var PATH = "/hello"

func hello(w http.ResponseWriter, req *http.Request) {
	helloName := "world"
	version := "SET_TO_APP_VERSION"

	repeatStr := os.Getenv("REPEAT")
	repeat := 1
	verboseStr := os.Getenv("VERBOSE")
	verbose := false

	if verboseStr == "true" || verboseStr == "TRUE" {
		verbose = true
	}

	if repeatStr != "" {
		var err error
		repeat, err = strconv.Atoi(repeatStr)
		if err != nil {
			repeat = 1
			fmt.Printf("Error: invalid value passed to repeat. Error from Atoi: %v", err)
			fmt.Println("Using value of 1.")
		}
	}

	for i := 0; i < repeat; i++ {
		fmt.Fprintf(w, "Hello %s\n", helloName)
	}

	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Values:")
	fmt.Fprintf(w, "  - Version: %s\n", version)
	fmt.Fprintf(w, "  - Repeat: %d\n", repeat)
	fmt.Fprintf(w, "  - Verbose: %t\n", verbose)

	if verbose {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "Verbose info:")
		fmt.Fprintln(w, "  - Version is hard coded into the image. When image is operator controlled, the correct version (and so image tag) to deploy comes from spec.repeat")
		fmt.Fprintln(w, "  - Repeat is an env, controlling how many times to say hello. When image is operator controlled it comes from spec.repeat")
		fmt.Fprintln(w, "  - Verbose is an env, controlling whether to display this extra explanatory text. When image is operator controlled it comes from spec.verbose")
	}
}

func main() {
	fmt.Printf("Running hello server on %s:%s\n", PATH, PORT)
	http.HandleFunc(PATH, hello)
	http.ListenAndServe(":"+PORT, nil)
}
