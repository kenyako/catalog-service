package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/kenyako/catalog-service/internal/app/config/section"
)

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

var Root Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file does not exist")
	}

	err = envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}
}
