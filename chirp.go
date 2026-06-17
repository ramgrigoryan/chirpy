package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ramgrigoryan/chirpy/internal/database"
)

type ChirpRequest struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}
type ChirpResponse struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfaneWords map[string]struct{}

func (cfg *apiConfig) CreateChirp(writer http.ResponseWriter, req *http.Request) {
	validatedChirpReq, err := ValidateChirp(writer, req)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to create chirp", err)
	}

	newChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   validatedChirpReq.Body,
		UserID: validatedChirpReq.UserID,
	})
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to save chirp", err)
		return
	}

	respondWithJSON(writer, http.StatusCreated, ChirpResponse{
		ID:        newChirp.ID,
		Body:      newChirp.Body,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		UserID:    newChirp.UserID,
	})
}

func ValidateChirp(writer http.ResponseWriter, req *http.Request) (ChirpRequest, error) {
	writer.Header().Add("content-type", "application/json")

	resp := ChirpRequest{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&resp)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "Couldn't decode params", err)
		return ChirpRequest{}, err
	}

	if len(resp.Body) > 140 {
		respondWithErr(writer, http.StatusBadRequest, "Chirp is too long", err)
		return ChirpRequest{}, err
	}

	// out, err := json.Marshal(cleanedData)
	// if err != nil {
	// 	respondWithErr(writer, http.StatusInternalServerError, "Error marshalling JSON", err)
	// 	return
	// }

	// writer.Write([]byte(out))

	return ChirpRequest{
		Body:   replaceProfaneWords(resp.Body),
		UserID: resp.UserID,
	}, nil
}

func replaceProfaneWords(text string) string {
	prohibitedWords := ProfaneWords{
		"kerfuffle": struct{}{},
		"sharbert":  struct{}{},
		"fornax":    struct{}{},
	}

	words := strings.Split(text, " ")
	for i, word := range words {
		if _, ok := prohibitedWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) GetChirps(writer http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "unable to fetch chirps", err)
	}

	resChirps := make([]ChirpResponse, len(dbChirps))
	for i, dbChirp := range dbChirps {
		resChirps[i] = ChirpResponse{
			ID:        dbChirp.ID,
			Body:      dbChirp.Body,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.ID,
		}
	}
	respondWithJSON(writer, http.StatusOK, resChirps)
}

func (cfg *apiConfig) GetChirp(writer http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("id")
	if chirpId == "" {
		respondWithErr(writer, http.StatusBadRequest, "chirp not specified", errors.New("empty chirp id"))
	}

	dbChirp, err := cfg.dbQueries.GetChirp(req.Context(), chirpId)
	if err != nil {
		respondWithErr(writer, http.StatusNotFound, "chirp not found", err)
	}

	respondWithJSON(writer, http.StatusOK, ChirpResponse{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
	})
}
