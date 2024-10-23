package utils

import (
	"errors"
	"math/rand"
)

const (
	lengthShortURL = 5
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// EncodeURL - кодируем URL.
func EncodeURL(url string) (string, error) {

	b := make([]byte, lengthShortURL)
	if url != "" {

		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		return string(b), nil
	}

	return "", errors.New("URL is empty")
}
