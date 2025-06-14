package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	STAGE        string
	RPC          map[string]string
	SMTPID       string
	SMTPPassword string
	DBUser       string
	DBPassword   string
	DBName       string
	DBHost       string
	DBPort       string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	return &Config{
		STAGE: os.Getenv("STAGE"),
		RPC: map[string]string{
			"Ethereum":         os.Getenv("ETH_RPC"),
			"Biance":           os.Getenv("BSC_RPC"),
			"GiantMammoth":     os.Getenv("GMMT_RPC"),
			"TestGiantMammoth": os.Getenv("TEST_GMMT_RPC"),
		},
		SMTPID:       os.Getenv("SMTP_ID"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
	}
}
