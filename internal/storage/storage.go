package storage

import model "github.com/IgorGreusunset/shortener/internal/app"

var Storage []model.URL = []model.URL{}

func WriteToStorage(record model.URL){
	Storage = append(Storage, record)
}

func ReadFromStorage(id string) *model.URL{
	var fullURL model.URL

	for _, v := range Storage {
		if v.ID == id {
			fullURL = v
		}
	}

	return &fullURL
}