package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	type errorReturn struct {
		Error string `json:"error"`
	}
	respondWithJSON(writer, code, errorReturn{Error: msg})
}

func respondWithJSON(writer http.ResponseWriter, code int, payload interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	log.Printf("Payload: %v", payload)
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(500)
		writer.Write([]byte(`{"error": "Error marshaling JSON response"}`))
		return
	}
	writer.WriteHeader(code)
	writer.Write(data)
}
