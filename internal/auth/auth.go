package auth

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(
		password,
		argon2id.DefaultParams,
	)
	if err != nil {
		fmt.Printf("error hashing password: %s\n", err)
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	checkHash, _, err := argon2id.CheckHash(password, hash)
	if err != nil {
		fmt.Printf("error checking password: %s\n", err)
		return false, err
	}
	return checkHash, nil
}
