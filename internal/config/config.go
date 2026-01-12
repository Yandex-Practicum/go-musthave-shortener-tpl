package config

import (
	"github.com/spf13/pflag"
)

type WebConfig struct {
	HostPort string
	BaseUrl  string
}

type Config struct {
	WebConfig WebConfig
}

func InitFlagConfig() *Config {
	cfg := &Config{}

	pflag.StringVar(&cfg.WebConfig.HostPort, "a", "localhost:8080", "server host")
	pflag.StringVar(&cfg.WebConfig.BaseUrl, "b", "http://localhost:8080", "base url")
	return cfg
}
