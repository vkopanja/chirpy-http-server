package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("test123")
	if err != nil {
		t.Fatalf("failed hashing password: %s", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	expected := "test123"
	hashed, err := HashPassword(expected)
	if err != nil {
		t.Fatalf("failed hashing password: %s", err)
	}
	check, err := CheckPasswordHash(expected, hashed)
	if err != nil {
		t.Fatalf("failed checking password: %s", err)
	}
	if !check {
		t.Fatalf("passwords do not match")
	}
}

func TestMakeJWT(t *testing.T) {
	_, err := MakeJWT(uuid.New(), "test", 10)
	if err != nil {
		t.Fatalf("failed creating jwt: %s", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tkn, err := MakeJWT(userID, "test", 10*time.Minute)
	if err != nil {
		t.Fatalf("failed creating jwt: %s", err)
	}

	tokenUserID, err := ValidateJWT(tkn, "test")
	if err != nil {
		t.Fatalf("failed validating jwt: %s", err)
	}

	if tokenUserID == uuid.Nil {
		t.Fatalf("failed validating jwt: %s", err)
	}

	if tokenUserID != userID {
		t.Fatalf("failed validating jwt, wrong subject extracted: %s", err)
	}
}
