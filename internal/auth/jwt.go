package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().In(time.UTC)),
		ExpiresAt: jwt.NewNumericDate(time.Now().In(time.UTC).Add(expiresIn)),
		Subject:   userID.String(),
	})
	tokenStr, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return "", err
	}
	return tokenStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("Error validating token: %v", err)
		return uuid.UUID{}, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		return uuid.UUID{}, err
	}
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		return uuid.UUID{}, err
	}
	return uuidUserID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", fmt.Errorf("Error getting Authorization Header Field")
	}
	prefix := "Bearer "
	if strings.HasPrefix(authString, prefix) {
		return strings.TrimPrefix(authString, prefix), nil
	}
	return "", fmt.Errorf("Authorization Field doesn't contain 'Bearer ' prefix")
}
