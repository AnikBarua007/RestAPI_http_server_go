package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	nt := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(nt),
		ExpiresAt: jwt.NewNumericDate(nt.Add(expiresIn)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	tokenClaims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}
	userID, err := uuid.Parse(tokenClaims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeader := strings.TrimSpace(headers.Get("Authorization"))
	if authorizationHeader == "" {
		return "", errors.New("missing authorization header")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, bearerPrefix))
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}
func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
func GetApiKey(headers http.Header) (string, error) {
	authorizationHeader := strings.TrimSpace(headers.Get("Authorization"))
	if authorizationHeader == "" {
		return "", errors.New("missing authorization header")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authorizationHeader, "ApiKey"))
	if token == "" {
		return "", errors.New("missing bearer token")
	}
	return token, nil
}
