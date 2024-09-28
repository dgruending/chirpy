package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func chirpValidationHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("Validation Request")
	type parameters struct {
		Chirp string `json:"body"`
	}
	// All responses are json
	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error while decoding parameters.")
		return
	}
	validateLength(writer, params.Chirp)
}

func validateLength(writer http.ResponseWriter, chirp string) {
	if len(chirp) <= 140 {
		cleanBadWords(writer, chirp)
	} else {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
	}
}

func cleanBadWords(writer http.ResponseWriter, chirp string) {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(chirp, " ")
	for ind, word := range chirpWords {
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpWords[ind] = "****"
		}
	}
	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respondWithJSON(writer, http.StatusOK, returnVal{CleanedBody: strings.Join(chirpWords, " ")})
}
