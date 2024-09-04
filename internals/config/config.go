package config

import (
	"log"
	"os"
)

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
}

func GetConfig() *Config {
	return &Config{
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD"),
		POSTGRES_USER:     getEnv("POSTGRES_USER"),
		POSTGRES_DB:       getEnv("POSTGRES_DB"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT"),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST"),
		PORT:              getEnv("PORT"),
	}
}
