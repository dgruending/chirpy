package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/dgruending/chirpy/internal/database"
	"github.com/google/uuid"
)

func validateLength(chirp string) bool {
	const chirpMaxLength = 140
	return len(chirp) < chirpMaxLength
}

func cleanBadWords(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(chirp, " ")
	for ind, word := range chirpWords {
		if slices.Contains(badWords, strings.ToLower(word)) {
			chirpWords[ind] = "****"
		}
	}
	return strings.Join(chirpWords, " ")
}

func (cfg *apiConfig) chirpHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Chirp  string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error while decoding Chirp.")
		return
	}

	if !validateLength(params.Chirp) {
		respondWithError(writer, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Chirp     string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	data, err := cfg.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{Body: params.Chirp, UserID: params.UserID})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating chirp")
		log.Printf("Error while creating user: %v", err)
		return
	}
	chirp := responseBody{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Chirp:     data.Body,
		UserID:    data.UserID,
	}
	respondWithJSON(writer, http.StatusCreated, chirp)
}
