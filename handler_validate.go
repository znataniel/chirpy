package main

import (
	"encoding/json"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type valid struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	if err := decoder.Decode(&ch); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "couldn't decode json parameters")
		return
	}

	if len(ch.Body) > 0 && len(ch.Body) <= 140 {
		res := valid{
			Valid: true,
		}
		respondJson(w, http.StatusOK, res)
		return
	}

	respondJsonError(w, http.StatusBadRequest, nil, "chirp is null or too long")
}
