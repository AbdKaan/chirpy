package auth

import "golang.org/x/crypto/bcrypt"

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
