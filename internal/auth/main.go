package auth

import (
	"errors"

	"github.com/alexedwards/argon2id"
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
