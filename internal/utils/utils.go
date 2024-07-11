package utils

import (
	"encoding/base64"
	"errors"
	"log"
)

func EncodeURL(url string) (string, error) {

	lenWord := 6
	if url != "" {
		var shortURL string
		startCoder := len(base64.StdEncoding.EncodeToString([]byte(url)))
		shortURL = base64.StdEncoding.EncodeToString([]byte(url))[startCoder-lenWord:]
		log.Println("URL encoded successfully")
		return shortURL, nil
	}

	return "", errors.New("URL is empty")
}
