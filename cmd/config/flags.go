package config

import (
	"flag"
	"os"
)

const (
	defaultServ = "localhost:8080"
	defaultBase =  "http://localhost:8080"
)

var (
	Serv string
	Base string
)

func ParseFlag() {

	Serv = os.Getenv("SERVER_ADDRESS")
	Base = os.Getenv("BASE_URL")
	

	servFlag := flag.String("a", defaultServ, "address  to run server")
	baseFlag := flag.String("b", defaultBase, "base address for short URL")
	flag.Parse()


	if Serv == "" {
		Serv = *servFlag
	}

	if Base == "" {
		Base = *baseFlag
	}

}