package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Database Database
}

func NewConfig() Config {
	return Config{
		App:      NewApp(),
		Database: NewDatabase(),
	}
}

type App struct {
	Port    string
	Timeout time.Duration
	GinMode string
}

func NewApp() App {
	timeoutStr := os.Getenv("APP_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal("Couldn't parse Timeout")
	}

	return App{
		Port:    os.Getenv("PORT"),
		Timeout: timeout,
		GinMode: os.Getenv("APP_GIN_MODE"),
	}
}

type Database struct {
	Name          string
	MigrationFile string
}

func NewDatabase() Database {
	return Database{
		Name:          os.Getenv("DB_NAME"),
		MigrationFile: os.Getenv("DB_MIGRATION"),
	}
}

type Token struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PublicKey            string
	PrivateKey           string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return NewConfig()
}
