package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type userEmail struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	e := userEmail{}

	if err := decoder.Decode(&e); err != nil {
		respondJsonError(w, http.StatusInternalServerError, err, "could not decode json data")
		return
	}

	createdUser, err := cfg.dbq.CreateUser(r.Context(), e.Email)
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
