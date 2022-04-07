package main

import (
	"fmt"
	"net/http"
	"os"
)

// PORT - server port
var PORT = "8080"

// PATH - server path
var PATH = "/hello"

func hello(w http.ResponseWriter, req *http.Request) {
	helloName := "world"
	envName := os.Getenv("HELLONAME")
	if envName != "" {
		helloName = envName
	}
	fmt.Fprintf(w, "Hello %s.\n\nBy the way, I'm the original 1.0 version of this awesome app.\n", helloName)
}

func main() {
	fmt.Printf("Running hello server on %s:%s\n", PATH, PORT)
	http.HandleFunc(PATH, hello)
	http.ListenAndServe(":"+PORT, nil)
}
