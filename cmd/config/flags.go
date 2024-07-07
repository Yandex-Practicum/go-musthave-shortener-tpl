package config

import (
	"flag"
)

type Config struct {
	Serv string
	Base string
}

func ParseFlag() *Config {
	ServAdd := flag.String("a", "localhost:8080", "address  to run server")
	BaseURL := flag.String("b", "http://localhost:8000", "base address for short URL")
	flag.Parse()

	return &Config{Serv: *ServAdd, Base: *BaseURL}
}