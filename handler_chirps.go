package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type chirpInput struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	ch := chirpInput{}
	if err := decoder.Decode(&ch); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not decode json data")
		return
	}

	if len(ch.Body) < 1 || len(ch.Body) > 140 {
		respondJsonError(w, http.StatusBadRequest, nil, "chirp is null or too long")
		return
	}

	createdChirp, err := cfg.dbq.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanChrip(ch.Body),
		UserID: ch.UserID,
	})
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not create chirp in db")
		return
	}

	jsonChirp := Chirp{
		ID:        createdChirp.ID,
		CreatedAt: createdChirp.CreatedAt,
		UpdatedAt: createdChirp.UpdatedAt,
		Body:      createdChirp.Body,
		UserID:    createdChirp.UserID,
	}

	respondJson(w, http.StatusCreated, jsonChirp)

}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	var jsonChirps []Chirp

	gotChirps, err := cfg.dbq.GetAllChirps(r.Context())
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not retrieve chirps from db")
		return
	}

	for _, c := range gotChirps {
		jsonChirps = append(jsonChirps, Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		})
	}

	respondJson(w, http.StatusOK, jsonChirps)

}
