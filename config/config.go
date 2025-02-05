package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	// JWT settings
	JWTSecret []byte

	// Database settings
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
)

// LoadConfig loads configuration from environment variables
func LoadConfig() error {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	// JWT settings
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}
	fmt.Printf("Loaded JWT secret (first 10 chars): %s...\n", jwtSecret[:10])
	JWTSecret = []byte(jwtSecret)

	// Database settings
	DBUser = getEnvWithDefault("DB_USER", "postgres")
	DBPassword = getEnvWithDefault("DB_PASSWORD", "0000")
	DBName = getEnvWithDefault("DB_NAME", "cbcexams")
	DBHost = getEnvWithDefault("DB_HOST", "localhost")
	DBPort = getEnvWithDefault("DB_PORT", "5432")

	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
