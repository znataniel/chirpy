package main

import "net/http"

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondJsonError(w, http.StatusForbidden, nil, "access denied")
		return
	}

	cfg.fileserverHits.And(0)

	if err := cfg.dbq.DeleteAllUsers(r.Context()); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not reset users table")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reset was succesful"))
}
