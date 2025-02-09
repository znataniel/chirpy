package main

import (
	"encoding/json"
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
	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(template, cfg.fileserverHits.Load())))
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

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type notValid struct {
		Error string `json:"error"`
	}

	type valid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	if err := decoder.Decode(&ch); err != nil {
		log.Printf("error decoding request body: %s", err)
		w.WriteHeader(500)
		return
	}

	var data []byte
	var err error
	w.Header().Set("Content-Type", "application/json")

	if len(ch.Body) > 0 && len(ch.Body) <= 140 {
		res := valid{
			Valid: true,
		}
		data, err = json.Marshal(res)
		if err != nil {
			log.Printf("could not encode data: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	} else {
		res := notValid{
			Error: "Chirp is too long",
		}
		data, err = json.Marshal(res)
		if err != nil {
			log.Printf("could not encode data: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
	}

	w.Write(data)
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
