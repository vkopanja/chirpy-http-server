package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
}

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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "chirpy",
			Subject:   userID.String(),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expiresIn)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		},
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	withClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(withClaims.Claims.(*Claims).Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	if headers.Get("Authorization") == "" {
		return "", fmt.Errorf("no authorization header")
	}

	bearerToken := headers.Get("Authorization")
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header")
	}

	return bearerToken[7:], nil
}

func MakeRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(randBytes), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	if headers.Get("Authorization") == "" {
		return "", fmt.Errorf("no authorization header")
	}

	bearerToken := headers.Get("Authorization")
	if !strings.HasPrefix(bearerToken, "ApiKey ") {
		return "", fmt.Errorf("invalid authorization header")
	}

	return bearerToken[7:], nil
}
