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
	var PortStr string
	var ginMode string

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
		fmt.Printf("Using local Google Application Credentials: %s\n", gacPath)

		PortStr = getEnvOrDefault("PORT", "8080")
		ginMode = getEnvOrDefault("GIN_MODE", "debug")

	} else {
		PortStr = os.Getenv("PORT")
		ginMode = "release"
	}

	if PortStr == "" {
		return nil, fmt.Errorf("PORT environment variable is not set")
	}

	readTimeout, _ := time.ParseDuration(getEnvOrDefault("READ_TIMEOUT", "15m"))
	writeTimeout, _ := time.ParseDuration(getEnvOrDefault("WRITE_TIMEOUT", "60s"))
	idleTimeout, _ := time.ParseDuration(getEnvOrDefault("IDLE_TIMEOUT", "5m"))

	projectID := getEnvOrDefault("PROJECT_ID", "")
	if projectID == "" {
		return nil, fmt.Errorf("PROJECT_ID environment variable is not set")
	}
	region := getEnvOrDefault("REGION", "")

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable is not set")
	}

	return &Config{
		ProjectID:  projectID,
		Region:     region,
		BucketName: bucketName,
		Server: ServerConfig{
			Port:         PortStr,
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
