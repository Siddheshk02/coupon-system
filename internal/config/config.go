package config

import (
	"os"
)

type Config struct {
	Port  string
	DBUrl string
}

var AppConfig Config

func Load() {
	// err := godotenv.Load("../../.env")
	// if err != nil {
	// 	log.Printf("Error loading .env file: %v", err)
	// }

	AppConfig = Config{
		Port:  os.Getenv("PORT"),
		DBUrl: os.Getenv("DB_URL"),
	}
}
