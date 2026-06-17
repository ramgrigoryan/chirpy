package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ramgrigoryan/chirpy/internal/auth"
	"github.com/ramgrigoryan/chirpy/internal/database"
)

type parameter struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (a *apiConfig) CreateUser(writer http.ResponseWriter, req *http.Request) {
	params := parameter{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to decode params", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to hash password", err)
		return
	}

	user, err := a.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to create user", err)
		return
	}

	respondWithJSON(writer, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
