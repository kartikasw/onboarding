package config

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Database Database
	Redis    Redis
	SMTP     SMTP
	Token    Token
}

func NewConfig() Config {
	return Config{
		App:      NewApp(),
		Database: NewDatabase(),
		Redis:    NewRedis(),
		SMTP:     NewSMTP(),
		Token:    NewToken(),
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
	Name         string
	Host         string
	Port         string
	Password     string
	User         string
	Timezone     string
	SslMode      string
	MigrationURL string
}

func NewDatabase() Database {
	return Database{
		Name:         os.Getenv("DB_NAME"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Password:     os.Getenv("DB_PASSWORD"),
		User:         os.Getenv("DB_USER"),
		Timezone:     os.Getenv("DB_TIMEZONE"),
		SslMode:      os.Getenv("DB_SSLMODE"),
		MigrationURL: os.Getenv("DB_MIGRATION_URL"),
	}
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedis() Redis {
	dbNum, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	return Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	}
}

type SMTP struct {
	Host     string
	Port     string
	FromName string
	Username string
	Password string
}

func NewSMTP() SMTP {
	return SMTP{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		FromName: os.Getenv("SMTP_FROM_NAME"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

type Token struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PublicKey            string
	PrivateKey           string
}

func NewToken() Token {
	encodedPublicKey := os.Getenv("TOKEN_PUBLIC_KEY")

	publicKey, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		log.Fatal("Couldn't encode public key")
	}

	encodedPrivateKey := os.Getenv("TOKEN_PRIVATE_KEY")
	privateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		log.Fatal("Couldn't encode private key")
	}

	accessDurationStr := os.Getenv("TOKEN_ACCESS_TOKEN_DURATION")
	accessDuration, err := time.ParseDuration(accessDurationStr)
	if err != nil {
		log.Fatal("Couldn't parse Access Token Duration")
	}

	refreshDurationStr := os.Getenv("TOKEN_REFRESH_TOKEN_DURATION")
	refreshDuration, err := time.ParseDuration(refreshDurationStr)
	if err != nil {
		log.Fatal("Couldn't parse Access Token Duration")
	}

	return Token{
		AccessTokenDuration:  accessDuration,
		RefreshTokenDuration: refreshDuration,
		PublicKey:            string(publicKey),
		PrivateKey:           string(privateKey),
	}
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return NewConfig()
}
