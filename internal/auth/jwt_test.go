package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTNoError(t *testing.T) {
	userID := uuid.New()
	_, err := MakeJWT(userID, "This is my secret", time.Duration(1000000))
	if err != nil {
		t.Fatalf("JWT Generation failed with: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "A secret key"
	tokenString, err := MakeJWT(userID, tokenSecret, time.Duration(10000000000000))
	if err != nil {
		t.Fatalf("JWT Generation failed with: %v", err)
	}
	tokenUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("JWT Validation failed with: %v", err)
	}
	if tokenUserID != userID {
		t.Fatalf("User ID Mismatch")
	}
}

func TestMakeJWTNoSecret(t *testing.T) {
	userID := uuid.New()
	_, err := MakeJWT(userID, string([]byte{}), time.Duration(10))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestMakeJWTNullDuration(t *testing.T) {
	userID := uuid.New()
	_, err := MakeJWT(userID, string([]byte{}), time.Duration(0))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestMakeJWTNegativDuration(t *testing.T) {
	userID := uuid.New()
	_, err := MakeJWT(userID, string([]byte{}), time.Duration(-10))
	if err != nil {
		t.Fatalf("%v", err)
	}
}

// can't figure out, how to make my wrapper function 'MakeJWT' fail

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "Correct Secret"
	wrongSecret := "Another secrect"
	tokenString, err := MakeJWT(userID, tokenSecret, time.Duration(1000000000000000000))
	if err != nil {
		t.Fatalf("JWT Generation failed with: %v", err)
	}
	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Fatalf("JWT Validation successfull with wrong secret")
	}
}

func TestValidateJWTExpiredDuration(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "Correct Secret"
	tokenString, err := MakeJWT(userID, tokenSecret, time.Duration(1))
	if err != nil {
		t.Fatalf("JWT Generation failed with: %v", err)
	}
	time.Sleep(2)
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatalf("JWT Validation successfull with expired duration")
	}
}

func TestGetBearerTokenCorrect(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer TOKEN_STRING")
	tokenString, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("Error getting Authorization Bearer Token")
	}
	if tokenString != "TOKEN_STRING" {
		t.Fatalf("Wrong token string: '%s' != 'TOKEN_STRING'", tokenString)
	}
}

func TestGetBearerTokenNoAuthField(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("No Error getting Authorization Bearer Token")
	}
}
func TestGetBearerTokenMalformedAuthFieldNoBearer(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	header.Add("Authorization", " TOKEN_STRING")
	if err == nil {
		t.Fatalf("No Error getting Authorization Bearer Token")
	}
}
func TestGetBearerTokenMalformedAuthFieldNoSpace(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	header.Add("Authorization", "BearerTOKEN_STRING")
	if err == nil {
		t.Fatalf("No Error getting Authorization Bearer Token")
	}
}
func TestGetBearerTokenMalformedAuthFieldWrongBearer(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	header.Add("Authorization", "Berer TOKEN_STRING")
	if err == nil {
		t.Fatalf("No Error getting Authorization Bearer Token")
	}
}
