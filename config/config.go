package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println("Successfully loaded .env file")
	}
}
