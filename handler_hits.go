package main

import (
	"fmt"
	"net/http"
)

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
