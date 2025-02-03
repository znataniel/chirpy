package main

import (
	"log"
	"net/http"
)

func getHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {

	const PORT = "8080"
	const ROOT = "static/"

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(ROOT))))
	serveMux.HandleFunc("/healthz", getHealthz)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: serveMux,
	}

	log.Printf("serving directory %s in port %s", ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
