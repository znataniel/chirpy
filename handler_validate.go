package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type valid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirp{}
	if err := decoder.Decode(&ch); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "couldn't decode json parameters")
		return
	}

	if len(ch.Body) > 0 && len(ch.Body) <= 140 {
		res := valid{
			CleanedBody: cleanChrip(ch.Body),
		}
		respondJson(w, http.StatusOK, res)
		return
	}

	respondJsonError(w, http.StatusBadRequest, nil, "chirp is null or too long")
}

func cleanChrip(body string) string {
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(body, " ")
	for i, w := range words {
		if slices.Contains(bannedWords, strings.ToLower(w)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")

}
