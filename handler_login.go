package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/znataniel/chirpy/internal/auth"
	"github.com/znataniel/chirpy/internal/database"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := userInput{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not decode json data")
		return
	}

	gotUser, err := cfg.dbq.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "incorrect email or password")
		return
	}

	passValid := auth.CheckPasswordHash(input.Password, gotUser.HashedPassword)
	if passValid != nil {
		respondJsonError(w, http.StatusUnauthorized, err, "incorrect email or password")
		return
	}

	genToken, err := auth.MakeJWT(gotUser.ID, cfg.secret)

	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.dbq.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: time.Now().AddDate(60, 0, 0).UTC(),
		UserID:    gotUser.ID,
	})
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not store refresh token")
		return
	}

	respondJson(w, http.StatusOK, User{
		ID:           gotUser.ID,
		CreatedAt:    gotUser.CreatedAt,
		UpdatedAt:    gotUser.UpdatedAt,
		Email:        gotUser.Email,
		IsChirpyRed:  gotUser.IsChirpyRed,
		Token:        genToken,
		RefreshToken: refreshToken,
	})
}
