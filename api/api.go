package api

import (
	"log"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("api.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
