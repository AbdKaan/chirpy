package auth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	testPw1 := "abcdefg1234lol"
	hashedPw1, err := HashPassword(testPw1)
	print(hashedPw1)
	if err != nil {
		t.Errorf("trying to hash: %v", err)
	}

	if result := CheckPasswordHash(testPw1,hashedPw1); result != nil {
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
