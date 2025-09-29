package auth

import (
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	params := argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	hashPass, err := argon2id.CreateHash(password, &params)
	if err != nil {
		return "", err
	}
	return hashPass, nil
}
func CheckPasswordHash(password, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, errors.New("password doesnt match")
	}
	return true, nil
}
func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	time := time.Now()
	expireTime := time.Add(expiresIn)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{time},
		ExpiresAt: &jwt.NumericDate{expireTime},
		Subject:   userId.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte{})
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	jwt.ParseWithClaims(tokenString)
}
