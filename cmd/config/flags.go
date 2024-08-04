package config

import (
	"flag"
	"os"
)

const (
	defaultServ = "localhost:8080"
	defaultBase =  "http://localhost:8080"
	defaultFile = "./short_url.json"
)

var (
	Serv string
	Base string
	File string
	DataBase string
)

func ParseFlag() {

	Serv = os.Getenv("SERVER_ADDRESS")
	Base = os.Getenv("BASE_URL")
	File = os.Getenv("FILE_STORAGE_PATH")
	DataBase = os.Getenv("DATABASE_DSN")
	

	servFlag := flag.String("a", defaultServ, "address  to run server")
	baseFlag := flag.String("b", defaultBase, "base address for short URL")
	fileFlag := flag.String("f", defaultFile, "path to file to save short urls")
	flag.StringVar(&DataBase, "d", "", "string for database connection")
	flag.Parse()

	//Проверяем наличие адресов в переменном окружении, если их нет - берем адреса из флагов.
	if Serv == "" {
		Serv = *servFlag
	}

	if Base == "" {
		Base = *baseFlag
	}

	if File == "" {
		File = *fileFlag
	}
}