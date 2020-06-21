package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateRandomBytes ...
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString ...
func GenerateRandomString(n int) string {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}
