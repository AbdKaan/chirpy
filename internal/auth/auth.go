package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPw), nil
}

func CheckPasswordHash(password, hashPw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPw), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		err = fmt.Errorf("signing token: %v", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	// handler needs to respond with 401 Unauthorized if there is an error
	if err != nil {
		err = fmt.Errorf("parsing claim: %v", err)
		return uuid.Nil, err
	}

	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.Claims)
	userIdStr, err := tokenClaim.Claims.GetSubject()
	if err != nil {
		err = fmt.Errorf("getting subject: %v", err)
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		err = fmt.Errorf("parsing user id: %v", err)
		return uuid.Nil, err
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	headerAuth := headers.Get("Authorization")
	if len(headerAuth) <= 7 {
		err := fmt.Errorf("authorization header format is wrong. Authorization header: %v", headerAuth)
		return "", err
	} 

	if headerAuth[:7] != "Bearer " {
		err := fmt.Errorf("authorization header format is wrong. Authorization header: %v", headerAuth)
		return "", err
	}

	return headerAuth[7:], nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	headerAuth := headers.Get("Authorization")
	if len(headerAuth) <= 7 {
		err := fmt.Errorf("authorization header format is wrong. Authorization header: %v", headerAuth)
		return "", err
	} 

	if headerAuth[:7] != "ApiKey " {
		err := fmt.Errorf("authorization header format is wrong. Authorization header: %v", headerAuth)
		return "", err
	}

	return headerAuth[7:], nil
}