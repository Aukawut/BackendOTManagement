package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadDatabaseConfig() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return ""
	}

	// Set up the connection string for SQL Server
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s&encrypt=disable&connection+timeout=30",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_NAME"))

	return connString
}
