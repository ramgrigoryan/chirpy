package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ramgrigoryan/chirpy/internal/auth"
)

type authParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) AuthHandler(writer http.ResponseWriter, req *http.Request) {
	authParams := authParams{}

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&authParams); err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to decode params", err)
		return
	}

	dbUser, err := cfg.dbQueries.GetUser(req.Context(), authParams.Email)
	if err != nil {
		respondWithErr(writer, http.StatusUnauthorized, "Incorrect email or password", errors.New("incorrect email or password"))
		return
	}

	match, err := auth.CheckPasswordHash(authParams.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to verify password", err)
		return
	}
	if !match {
		respondWithErr(writer, http.StatusUnauthorized, "Incorrect email or password", errors.New("incorrect email or password"))
		return
	}

	respondWithJSON(writer, http.StatusOK, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	})
}
