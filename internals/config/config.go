package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var GlobalConfig *Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	GlobalConfig = GetConfig()
}
func getEnv(env string) string {
	data, ok := os.LookupEnv(env)
	if !ok {
		log.Fatalf("%s missing from environments", env)
	}
	return data
}

type Config struct {
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_DB       string
	POSTGRES_PORT     string
	POSTGRES_HOST     string
	PORT              string
	CONNECTION_STR    string
	ENV               string
	SMTP_USERNAME     string
	SMTP_PASSWORD     string
	SMTP_HOST         string
	SMTP_PORT         string
	SMTP_EMAIL        string
	FRONTEND_URL      string
}

func GetConfig() *Config {
	return &Config{
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     getEnv("POSTGRES_USER"),
		POSTGRES_DB:       getEnv("POSTGRES_DB"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT"),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST"),
		PORT:              getEnv("PORT"),
		CONNECTION_STR:    getEnv("CONNECTION_STR"),
		ENV:               getEnv("ENV"),
		SMTP_USERNAME:     getEnv("SMTP_USERNAME"),
		SMTP_PASSWORD:     getEnv("SMTP_PASSWORD"),
		SMTP_HOST:         getEnv("SMTP_HOST"),
		SMTP_PORT:         getEnv("SMTP_PORT"),
		SMTP_EMAIL:        getEnv("SMTP_EMAIL"),
		FRONTEND_URL:      getEnv("FRONTEND_URL"),
	}
}
