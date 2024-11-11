package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	testPw1 := "abcdefg1234lol"
	hashedPw1, err := HashPassword(testPw1)
	if err != nil {
		t.Errorf("trying to hash: %v", err)
	}

	if result := CheckPasswordHash(testPw1, hashedPw1); result != nil {
		t.Errorf("hash and pw not equal: %v", result)
	}

	testPw2 := "abc"
	hashedPw2, err := HashPassword(testPw2)
	if err != nil {
		t.Errorf("trying to hash: %v", err)
	}

	if result := CheckPasswordHash(testPw2, hashedPw2); result != nil {
		t.Errorf("hash and pw not equal: %v", result)
	}
}

func TestJWT(t *testing.T) {
	id := uuid.New()
	tokenSecret := "verysecret"
	tokenStr, err := MakeJWT(id, tokenSecret, time.Second*3)
	if err != nil {
		t.Errorf("trying to make jwt: %v", err)
	}

	// validate
	validatedUuid, err := ValidateJWT(tokenStr, tokenSecret)
	if err != nil {
		t.Errorf("validating jwt: %v", err)
	}
	// make sure validated id and original id is same
	if validatedUuid != id {
		t.Errorf("validated id is not equal to the original")
	}

	// shouldn't be validated with wrong secret key
	_, err = ValidateJWT(tokenStr, "fakenews")
	if err == nil {
		t.Errorf("error is nil with wrong token secret: %v", err)
	}

	time.Sleep(3 * time.Second)
	// token should be expired
	_, err = ValidateJWT(tokenStr, tokenSecret)
	if err == nil {
		t.Errorf("token should be expired but error is nil: %v", err)
	}
}

func TestGetBearerToken(t *testing.T) {
	token := "Bearer mytoken"
	header := http.Header{"Authorization": {token}}
	bearerToken, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("header not starting with 'Bearer ': %v", err)
	}

	if bearerToken != "mytoken" {
		t.Errorf("found tokens are not equal, %s != %s", bearerToken, token)
	}
}
