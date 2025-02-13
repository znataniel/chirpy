package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/znataniel/chirpy/internal/auth"
	"github.com/znataniel/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token" json:"-"`
	RefreshToken string    `json:"refresh_token" json:"-"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	e := userInput{}

	if err := decoder.Decode(&e); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not decode json data")
		return
	}

	pass, err := auth.HashPassword(e.Password)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not hash password")
		return
	}

	createdUser, err := cfg.dbq.CreateUser(r.Context(), database.CreateUserParams{
		Email:          e.Email,
		HashedPassword: pass,
	})
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not create user in db")
		return
	}

	u := User{
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	}
	respondJson(w, http.StatusCreated, u)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type updateUserInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	decoder := json.NewDecoder(r.Body)
	input := updateUserInput{}
	if err := decoder.Decode(&input); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not decode json data")
		return
	}

	hashedPass, err := auth.HashPassword(input.Password)
	if err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "error hashing password")
		return
	}

	updatedUser, err := cfg.dbq.UpdateUserAndPassword(
		r.Context(),
		database.UpdateUserAndPasswordParams{
			ID:             userID,
			Email:          input.Email,
			HashedPassword: hashedPass,
			UpdatedAt:      time.Now().UTC(),
		},
	)
	if err != nil {
		respondJsonError(
			w,
			http.StatusInternalServerError,
			err,
			"failed to update email & pass or user not found",
		)
	}

	respondJson(w, http.StatusOK, User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	})
}
