package config

import (
	"fmt"
	"log"
	"os"
	"strconv" // For parsing numbers/booleans from string env vars

	"github.com/joho/godotenv"
)

// Config holds all application configuration settings
type Config struct {
	DBUser       string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
	JWTSecretKey string
	AppPort      string
	MailMailer   string
	MailHost     string
	MailPort     string
	MailUsername string
	MailPassword string
	// Add any other configuration variables here
}

// LoadConfig loads environment variables from .env and populates the Config struct
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
		AppPort:      os.Getenv("APP_PORT"),
		MailMailer:   os.Getenv("MAIL_MAILER"),
		MailHost:     os.Getenv("MAIL_HOST"),
		MailPort:     os.Getenv("MAIL_PORT"),
		MailUsername: os.Getenv("MAIL_USERNAME"),
		MailPassword: os.Getenv("MAIL_PASSWORD"),
	}

	// Basic validation for critical config (you can add more comprehensive checks)
	if cfg.DBUser == "" || cfg.DBHost == "" ||
		cfg.DBPort == "" || cfg.DBName == "" || cfg.JWTSecretKey == "" ||
		cfg.AppPort == "" {
		log.Fatal("One or more critical environment variables are missing. Please check your .env file.")
	}

	// Example of parsing a boolean or integer if you had one in .env
	// debugModeStr := os.Getenv("DEBUG_MODE")
	// cfg.DebugMode, _ = strconv.ParseBool(debugModeStr)

	return cfg
}

// GetDBPortAsInt returns the DBPort as an integer, handling errors
func (c *Config) GetDBPortAsInt() (int, error) {
	port, err := strconv.Atoi(c.DBPort)
	if err != nil {
		return 0, fmt.Errorf("invalid DB_PORT in config: %v", err)
	}
	return port, nil
}