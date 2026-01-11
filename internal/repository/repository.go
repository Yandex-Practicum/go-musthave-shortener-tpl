package repository

import "github.com/google/uuid"

type Repository interface {
	Shorten(string) (string, error)
	Unshorten(string) (string, error)
}

var bd = make(map[string]string)
var ud = uuid.New()

func Shorten(url string) (string, error) {
	shurl := ud.String()
	bd[shurl] = url
	return shurl, nil
}

func Unshorten(shurl string) (string, error) {
	return bd[shurl], nil
}
