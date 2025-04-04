package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func ConfigLoad() *Config {
	if _, err := os.Stat("conf/local.yml"); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig("conf/local.yml", &cfg); err != nil {
		log.Fatalf("config file read error: %s", err)
	}

	return &cfg
}
