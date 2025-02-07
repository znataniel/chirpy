package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsFs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		h.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getFsHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) resetFsHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.And(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits reset was successful"))
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {

	const PORT = "8080"
	const ROOT = "static/"

	cfg := apiConfig{}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsFs(http.StripPrefix("/app", http.FileServer(http.Dir(ROOT)))))
	serveMux.HandleFunc("/healthz", getHealthz)
	serveMux.HandleFunc("/metrics", cfg.getFsHits)
	serveMux.HandleFunc("/reset", cfg.resetFsHits)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: serveMux,
	}

	log.Printf("serving directory %s in port %s", ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
