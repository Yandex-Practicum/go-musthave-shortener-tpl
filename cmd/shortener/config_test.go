package main

import (
	"testing"
)

// TestParse - тестирует парсинг конфигурационной строки.
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		addrServ string
		baseURL  string
		logLevel string
		pathFile string
		addrDB   string
	}{
		{
			name:     "successful",
			addrServ: ":8080",
			baseURL:  "http://localhost:8080",
			logLevel: "info",
			pathFile: "",
			addrDB:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfigs()
			cfg.AddrServer = tt.addrServ
			cfg.BaseURL = tt.baseURL
			cfg.LogLevel = tt.logLevel
			cfg.PathFile = tt.pathFile
			cfg.AddrDB = tt.addrDB
			cfg.Parse()

			if cfg.AddrServer != tt.addrServ {
				t.Errorf("Ожидали %v, пришли %v", tt.addrServ, cfg.AddrServer)
			}
			if cfg.BaseURL != tt.baseURL {
				t.Errorf("Ожидали %v, пришли %v", tt.baseURL, cfg.BaseURL)
			}
			if cfg.LogLevel != tt.logLevel {
				t.Errorf("Ожидали %v, пришли %v", tt.logLevel, cfg.LogLevel)
			}
			if cfg.PathFile != tt.pathFile {
				t.Errorf("Ожидали %v, пришли %v", tt.pathFile, cfg.PathFile)
			}
			if cfg.AddrDB != tt.addrDB {
			}
		})
	}
}
