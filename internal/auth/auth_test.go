package auth

import "testing"

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
