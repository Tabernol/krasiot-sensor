package oracledb

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type OracleConfig struct {
	User     string
	Password string
	ConnStr  string
	Wallet   string
	LibDir   string
}

func LoadOracleConfig() (OracleConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Using system env vars (no .env file found)")
	}
	return OracleConfig{
		User:     os.Getenv("ADB_USERNAME"),
		Password: os.Getenv("ADB_PASSWORD"),
		ConnStr:  os.Getenv("ADB_CONN_STRING"),     // full tcps://... string
		Wallet:   os.Getenv("ADB_WALLET_LOCATION"), // optional, maybe set in libDir
		LibDir:   os.Getenv("ADB_LIB_DIR"),
	}, nil
}
