package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/chirpy/internal/auth"
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
		Body string `json:"body"`
	}

	// jwt authorization
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "could not read token from header")
		return
	}
	tokenUserID, ok := auth.ValidateJWT(bearerToken, cfg.secret)
	if ok != nil {
		respondJsonError(w, http.StatusUnauthorized, ok, "token provided is not valid")
		return
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
		UserID: tokenUserID,
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

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	if s == "" {
		cfg.getAllChirps(w, r)
		return
	}

	userID, err := uuid.Parse(s)
	if err != nil {
		respondJsonError(w, http.StatusBadRequest, err, "invalid uuid in query parameter")
		return
	}
	cfg.getChirpsByUserID(w, r, userID)
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	gotChirps, err := cfg.dbq.GetAllChirps(r.Context())
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not retrieve chirps from db")
		return
	}

	var jsonChirps []Chirp
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

func (cfg *apiConfig) getChirpById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not read id value from path")
		return
	}

	ch, err := cfg.dbq.GetChirpById(r.Context(), id)
	if err != nil {
		respondJsonError(w, http.StatusNotFound, err, "chirp not found")
		return
	}

	respondJson(w, http.StatusOK, Chirp{
		ID:        ch.ID,
		CreatedAt: ch.CreatedAt,
		UpdatedAt: ch.UpdatedAt,
		Body:      ch.Body,
		UserID:    ch.UserID,
	})
}

func (cfg *apiConfig) getChirpsByUserID(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	gotChirps, err := cfg.dbq.GetAllChirpsById(r.Context(), userID)
	if err != nil {
		respondJsonError(w, http.StatusNotFound, err, "user not found")
		return
	}

	var jsonChirps []Chirp
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

func (cfg *apiConfig) deleteChirpById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not read id value from path")
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "authorization header not found")
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "could not validate access")
		return
	}

	ch, err := cfg.dbq.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondJsonError(w, http.StatusNotFound, err, "chirp was not found")
		return
	}
	if ch.UserID != userID {
		respondJsonError(w, http.StatusForbidden, err, "chirp does not belong to this user")
		return
	}

	if err := cfg.dbq.DeleteChirpByID(r.Context(), ch.ID); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "chirp could not be deleted")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
