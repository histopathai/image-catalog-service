package config

import (
	"fmt"
	"log"
	"os"
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
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	GINMode      string
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("ENV")

	if env == "LOCAL" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("❌ Failed to load .env file: %v", err)
		}

		gacPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if gacPath == "" {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
		}
		if _, err := os.Stat(gacPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS file does not exist at path: %s", gacPath)
		}
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gacPath)
		fmt.Printf("✅ Using local Google Application Credentials: %s\n", gacPath)
	}

	port := getEnvOrDefault("PORT", "8080")
	if port == "" {
		return nil, fmt.Errorf("PORT is required")
	}

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}

	region := os.Getenv("REGION")
	if region == "" {
		return nil, fmt.Errorf("REGION environment variable is not set")
	}

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable is not set")
	}

	readTimeout, _ := time.ParseDuration(getEnvOrDefault("READ_TIMEOUT", "15m"))
	writeTimeout, _ := time.ParseDuration(getEnvOrDefault("WRITE_TIMEOUT", "60s"))
	idleTimeout, _ := time.ParseDuration(getEnvOrDefault("IDLE_TIMEOUT", "5m"))
	ginMode := getEnvOrDefault("GIN_MODE", "release")

	return &Config{
		ProjectID:  projectID,
		Region:     region,
		BucketName: bucketName,
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			GINMode:      ginMode,
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
