package repository

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type DbConfig struct {
	Username   string
	Password   string
	BaseURL    string
	SchemaPath string
}

func LoadDBConfig() (*DbConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &DbConfig{
		Username:   os.Getenv("ADB_USERNAME"),
		Password:   os.Getenv("ADB_PASSWORD"),
		BaseURL:    os.Getenv("ADB_BASE_URL"),
		SchemaPath: os.Getenv("ADB_SCHEMA_PATH"),
	}, nil
}
