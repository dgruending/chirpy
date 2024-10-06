package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgruending/chirpy/internal/auth"
	"github.com/dgruending/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type requestParameters struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) createUserHandler(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	params := requestParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding request")
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Create User password error: ")
		respondWithError(writer, http.StatusBadRequest, "Password error")
		return
	}
	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.dbQueries.CreateUser(request.Context(), userParams)
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

func (cfg *apiConfig) loginHandler(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	params := requestParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error decoding request")
		return
	}
	user, err := cfg.dbQueries.GetUser(request.Context(), params.Email)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(writer, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > 60*60 {
		params.ExpiresInSeconds = 360
	}
	duration, err := time.ParseDuration(fmt.Sprintf("%ds", params.ExpiresInSeconds))
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error parsing Expiration duration")
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.serverSecret, duration)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error creating JWT Token")
		return
	}
	userPayload := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	respondWithJSON(writer, http.StatusOK, userPayload)
}
