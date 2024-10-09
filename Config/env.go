package config

import (
	"fmt"
	"os"
	"strconv"
	// "log"
	// "github.com/joho/godotenv"
)

var (
	// DB is the database connection string
	DBPort     string
	DBHost     string
	DBUser     string
	DBName     string
	DBPassword string
	APIID      int32
	APIHash    string

	PhoneNumber string
	BotToken    string

	GeminAPIKey string
)

func EnvInit() {

	DBPort = os.Getenv("DB_PORT")
	if DBPort == "" {
		fmt.Println("PORT is not set. Using default port 8080")
	}

	DBHost = os.Getenv("DB_HOST")
	if DBHost == "" {
		fmt.Println("HOST is not set. Using default host localhost")
	}

	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		fmt.Println("NAME is not set. Using default name postgres")
	}

	DBUser = os.Getenv("DB_USER")
	if DBUser == "" {
		fmt.Println("USER is not set. Using default user postgres")
	}

	DBPassword = os.Getenv("DB_PASSWORD")
	if DBPassword == "" {
		fmt.Println("PASSWORD is not set. Using default password postgres")
	}

	apiIDStr := os.Getenv("API_ID")
	if apiIDStr == "" {
		fmt.Println("API_ID is not set. Using default API_ID")
	} else {
		apiID, err := strconv.Atoi(apiIDStr)
		if err != nil {
			fmt.Println("Invalid API_ID. Using default API_ID")
		} else {
			APIID = int32(apiID)
		}
	}

	APIHash = os.Getenv("API_HASH")
	if APIHash == "" {
		fmt.Println("API_HASH is not set. Using default API_HASH")
	}

	PhoneNumber = os.Getenv("PHONE_NUMBER")
	if PhoneNumber == "" {
		fmt.Println("PHONE_NUMBER is not set. Using default phone number")
	}

	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		fmt.Println("BOT_TOKEN is not set. Using default bot token")
	}

	GeminAPIKey = os.Getenv("GEMINI_API_KEY")
	if GeminAPIKey == "" {
		fmt.Println("GEMIN_API_KEY is not set. Using default Gemin API key")
	}

}
