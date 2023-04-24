package config

import (
	"api/app/utils"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.ErrorLogger.Fatal("Error loading .env file")
	}
}
