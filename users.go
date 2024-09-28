package main

import (
	"encoding/json"
	"log"
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

func (cfg *apiConfig) createUserHandler(writer http.ResponseWriter, request *http.Request) {
	type emailBody struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(request.Body)
	emailParam := emailBody{}
	err := decoder.Decode(&emailParam)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding email")
		return
	}
	user, err := cfg.dbQueries.CreateUser(request.Context(), emailParam.Email)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating user")
		log.Printf("Error while creating user: %v", err)
		return
	}
	// map to own User struct for control over json keys
	new_user := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(writer, http.StatusCreated, new_user)
}
