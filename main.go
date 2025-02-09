package main

import (
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

func main() {

	const PORT = "8080"
	const ROOT = "static/"

	cfg := apiConfig{}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsFs(http.StripPrefix("/app", http.FileServer(http.Dir(ROOT)))))
	serveMux.HandleFunc("GET /api/healthz", getHealthz)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)
	serveMux.HandleFunc("GET /admin/metrics", cfg.getFsHits)
	serveMux.HandleFunc("POST /admin/reset", cfg.resetFsHits)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: serveMux,
	}

	log.Printf("serving directory %s in port %s", ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
