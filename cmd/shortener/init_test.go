package main

import (
	"github.com/kamencov/go-musthave-shortener-tpl/internal/service"
	"reflect"
	"testing"
)

// MockPstStorage - заглушка для db.NewPstStorage
type MockPstStorage struct {
	service.Storage
}

func TestInitDB(t *testing.T) {
	cases := []struct {
		name     string
		db       string
		file     string
		expected string
	}{
		{
			name:     "WithDatabase",
			db:       "test_db",
			file:     "",
			expected: "*db.PstStorage",
		},
		{
			name:     "WithFile",
			db:       "",
			file:     "test_file",
			expected: "*filestorage.SaveFile",
		},
		{
			name:     "WithMap",
			db:       "",
			file:     "",
			expected: "*mapstorage.MapStorage",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			storage := initDB(tt.db, tt.file)
			if reflect.TypeOf(storage).String() != tt.expected {
				t.Errorf("Expected storage to be *db.PstStorage, got %T", storage)
			}
		})
	}
}
