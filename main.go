package main

import (
	"log"
	"net/http"
)

func main() {

	const PORT = "8080"
	const ROOT = "."

	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(ROOT)))

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: serveMux,
	}

	log.Printf("serving directory %s in port %s", ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
