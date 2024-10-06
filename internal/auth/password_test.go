package auth

import "testing"

func TestHashPasswordCorrect(t *testing.T) {
	password := "Simple@12Test"
	_, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hash failed with: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "Simple@12Test"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hash failed with: %v", err)
	}
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Check failed with: %v", err)
	}
}

func TestTooLongPassword(t *testing.T) {
	password := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	_, err := HashPassword(password)
	if err == nil {
		t.Fatalf("Expected error got nil.")
	}
}

func TestWrongPassword(t *testing.T) {
	password := "correct_password"
	wrongPassword := "Password123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Hash failed with: %v", err)
	}
	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Fatalf("Wrong password for hash")
	}
}
