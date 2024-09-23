package config

import (
	"fmt"
	"os"
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

}
