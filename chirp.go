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

type chirpEntry struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Chirp     string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func translateToChirp(dbChirp database.Chirp) chirpEntry {
	return chirpEntry{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Chirp:     dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
}

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

func (cfg *apiConfig) createChirpHandler(writer http.ResponseWriter, request *http.Request) {
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

	data, err := cfg.dbQueries.CreateChirp(request.Context(), database.CreateChirpParams{Body: params.Chirp, UserID: params.UserID})
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating chirp")
		log.Printf("Error while creating user: %v", err)
		return
	}
	chirp := translateToChirp(data)
	respondWithJSON(writer, http.StatusCreated, chirp)
}

func (cfg *apiConfig) getChirpsHandler(writer http.ResponseWriter, request *http.Request) {
	dbChirps, err := cfg.dbQueries.GetChirps(request.Context())
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	chirps := []chirpEntry{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, translateToChirp(dbChirp))
	}
	respondWithJSON(writer, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("Single Chirp Request")
	chirpID, err := uuid.Parse(request.PathValue("chirpID"))
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "Error while parsing chirp ID")
		return
	}
	dbChirp, err := cfg.dbQueries.GetChirpByID(request.Context(), chirpID)
	if err != nil {
		respondWithError(writer, http.StatusNotFound, "Chirp not found")
		return
	}
	respondWithJSON(writer, http.StatusOK, translateToChirp(dbChirp))
}
