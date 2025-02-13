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
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}
	respondJson(w, http.StatusCreated, u)
}
