package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ProjectID  string
	Region     string
	BucketName string
	Server     ServerConfig
}

type ServerConfig struct {
	Port         int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	GINMode      string // Added GIN_MODE for Gin framework configuration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Failed to load .env file: %v", err)
	}

	env := os.Getenv("ENV")
	if env == "LOCAL" {
		gacPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if gacPath == "" {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
		}
		if _, err := os.Stat(gacPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS file does not exist at path: %s", gacPath)
		}
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gacPath)
		fmt.Printf("Using local Google Application Credentials: %s\n", gacPath)
	}

	readTimeout, _ := time.ParseDuration(getEnvOrDefault("READ_TIMEOUT", "15m"))
	writeTimeout, _ := time.ParseDuration(getEnvOrDefault("WRITE_TIMEOUT", "60s"))
	idleTimeout, _ := time.ParseDuration(getEnvOrDefault("IDLE_TIMEOUT", "5m"))

	projectID := getEnvOrDefault("GCP_PROJECT_ID", "")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT_ID environment variable is not set")
	}
	region := getEnvOrDefault("GCP_REGION", "")

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable is not set")
	}

	return &Config{
		ProjectID:  projectID,
		Region:     region,
		BucketName: bucketName,
		Server: ServerConfig{
			Port:         int(getEnvAsInt("PORT", 8080)),
			Environment:  getEnvOrDefault("ENV", "development"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			GINMode:      getEnvOrDefault("GIN_MODE", "debug"), // Default to debug mode
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		fmt.Printf("Error parsing environment variable %s: %v. Using default value: %d\n", key, err, defaultValue)
		return defaultValue
	}
	return value
}
