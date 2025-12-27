package config

import (

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUSer     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}

func NewENVConfig() *EnvConfig {
	err := godotenv.Load()

	if err != nil{
		log.Fatalf("unable to load .env: %e", err)
	}
	
	config := &EnvConfig{}
	
	if err := env.Parse(config); err != nil{
		log.Fatalf("unable to load variable .env: %e", err)
	}
	return config
}