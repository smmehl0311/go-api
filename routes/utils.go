package routes

import (
	"crypto/rand"
	"crypto/sha256"
)

func GetToken() ([]byte, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	hasher.Write(b)
	return hasher.Sum(nil), nil
}
