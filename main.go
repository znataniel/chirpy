package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/znataniel/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbq            *database.Queries
	platform       string
	secret         string
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

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	cfg := apiConfig{
		dbq:      database.New(db),
		platform: os.Getenv("PLATFORM"),
		secret:   os.Getenv("SECRET"),
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", cfg.middlewareMetricsFs(http.StripPrefix("/app", http.FileServer(http.Dir(ROOT)))))
	serveMux.HandleFunc("GET /api/healthz", getHealthz)
	serveMux.HandleFunc("GET /admin/metrics", cfg.getFsHits)
	serveMux.HandleFunc("POST /admin/reset", cfg.reset)
	serveMux.HandleFunc("POST /api/users", cfg.createUser)
	serveMux.HandleFunc("POST /api/chirps", cfg.createChirp)
	serveMux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpById)
	serveMux.HandleFunc("POST /api/login", cfg.login)
	serveMux.HandleFunc("POST /api/refresh", cfg.refresh)
	serveMux.HandleFunc("POST /api/revoke", cfg.revoke)

	server := http.Server{
		Addr:    ":" + PORT,
		Handler: serveMux,
	}

	log.Printf("serving directory %s in port %s", ROOT, PORT)
	log.Fatal(server.ListenAndServe())
}
