package config

import (
	"flag"
)

var (
	Serv string
	Base string
)

func ParseFlag() {
	flag.StringVar(&Serv, "a", "localhost:8080", "address  to run server")
	flag.StringVar(&Base, "b", "http://localhost:8080", "base address for short URL")
	flag.Parse()

	if Serv == "" {
		Serv = "localhost:8080"
	}

	if Base == "" {
		Base = "http://localhost:8080"
	}

}