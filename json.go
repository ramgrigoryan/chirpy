package main

import (
	"encoding/json"
	"net/http"

	"log"
)

type Err struct {
	Error string `json:"error"`
}

func respondWithErr(writer http.ResponseWriter, code int, message string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Respinding with 5XX code %s", message)
	}

	respondWithJSON(writer, code, Err{
		Error: message,
	})

}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	writer.Header().Set("content-type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(code)
	writer.Write([]byte(data))
}
