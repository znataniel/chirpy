package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/znataniel/chirpy/internal/auth"
	"github.com/znataniel/chirpy/internal/database"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	type upgradeUserInput struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	// authorization
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "no api key found")
		return
	}
	if apiKey != cfg.polka_key {
		respondJsonError(w, http.StatusUnauthorized, err, "invalid api key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	input := upgradeUserInput{}
	if err := decoder.Decode(&input); err != nil {
		respondJsonError(w, http.StatusBadRequest, err, "could not decode json body")
		return
	}

	if input.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userUUID, err := uuid.Parse(input.Data.UserID)
	if err != nil {
		respondJsonError(
			w,
			http.StatusInternalServerError,
			err,
			"could not parse user id",
		)
		return
	}

	if err := cfg.dbq.SetChirpyRedStatusById(
		r.Context(),
		database.SetChirpyRedStatusByIdParams{
			ID:          userUUID,
			IsChirpyRed: true,
		},
	); err != nil {
		respondJsonError(w, http.StatusNotFound, err, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
