package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const TokenTypeAccess TokenType = "chirpy-access"

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// HashPassword generates a bcrypt hash of the given password.
// It returns the hashed password as a string or an error if hashing fails.
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash compares a plaintext password with its bcrypt hashed version.
// It returns an error if the password does not match the hash.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT generates a JWT token for the given user ID with the provided secret and expiration duration.
// It returns the signed token string or an error if token creation fails.
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

// ValidateJWT validates a JWT token using the provided secret.
// It returns the user ID embedded in the token if valid, or an error if the token is invalid.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil
}

// GetAuthToken extracts the authentication token from the HTTP Authorization header using the specified scheme.
// It returns the token or an error if the header is missing or malformed.
func GetAuthToken(headers http.Header, scheme string) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != scheme {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// GetBearerToken extracts the Bearer token from the HTTP Authorization header.
// It returns the token string or an error if the token is missing or malformed.
func GetBearerToken(headers http.Header) (string, error) {
	return GetAuthToken(headers, "Bearer")
}

// GetAPIKey extracts the API key from the HTTP Authorization header using the "ApiKey" scheme.
// It returns the API key string or an error if not found.
func GetAPIKey(headers http.Header) (string, error) {
	return GetAuthToken(headers, "ApiKey")
}

// MakeRefreshToken generates a secure random 256-bit token encoded in hexadecimal.
// It returns the refresh token or an error if token generation fails.
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
