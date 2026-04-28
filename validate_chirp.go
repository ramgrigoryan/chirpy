package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request struct {
	Body string `json:"body"`
}

type CleanerdData struct {
	ClearnedBody string `json:"cleaned_body"`
}

type ProfaneWords map[string]struct{}

func ValidateChirp(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("content-type", "application/json")

	resp := Request{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&resp)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}

	if len(resp.Body) > 140 {
		respondWithErr(writer, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	cleanedData := CleanerdData{
		ClearnedBody: replaceProfaneWords(resp.Body),
	}

	out, err := json.Marshal(cleanedData)
	if err != nil {
		respondWithErr(writer, http.StatusInternalServerError, "Error marshalling JSON", err)
		return
	}

	writer.Write([]byte(out))

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
