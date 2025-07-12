package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server ServerConfig
	GCS    GCSConfig
	IPS    IPSConfig
}

type GCSConfig struct {
	OriginalBucket   string   //Cold storage bucket for original images
	SupportedFormats []string //Supported image formats
	MaxFileSize      int64    //Maximum file size for uploads
	BucketLocation   string   //Location of the GCS bucket
	StorageClass     string   //Storage class for the GCS bucket
}

type ServerConfig struct {
	PROJECT_ID   string
	Port         int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type IPSConfig struct { // Image Processing Server Config
	Host     string
	Port     int
	Username string
	Password string
}

func LoadConfig() (*Config, error) {

	env_location := os.Getenv("ENV_LOCATION")
	if env_location == "LOCAL" {
		//look at credential file
		gac_path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if gac_path == "" {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set")
		}
		if _, err := os.Stat(gac_path); os.IsNotExist(err) {
			return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS file does not exist at path: %s", gac_path)
		}
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gac_path)
	}

	// Parser Timeouts from environment variables or set defaults
	readTimeout, _ := time.ParseDuration(getEnvOrDefault("READ_TIMEOUT", "15m"))
	writeTimeout, _ := time.ParseDuration(getEnvOrDefault("WRITE_TIMEOUT", "60s"))
	idleTimeout, _ := time.ParseDuration(getEnvOrDefault("IDLE_TIMEOUT", "5m"))

	return &Config{
		Server: ServerConfig{
			Port:         int(getEnvAsInt("PORT", 8080)),
			Environment:  getEnvOrDefault("ENV", "development"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
			PROJECT_ID:   getEnvOrDefault("PROJECT_ID", "your-project-id"),
		},
		GCS: GCSConfig{
			OriginalBucket:   os.Getenv("GCS_ORIGINAL_BUCKET"),
			SupportedFormats: strings.Split(getEnvOrDefault("SUPPORTED_FORMATS", "jpg,jpeg,png,tiff,tif,svs,ndpi,mrxs"), ","),
			MaxFileSize:      getEnvAsInt("MAX_FILE_SIZE", 5368709120), // 5GB default
			BucketLocation:   getEnvOrDefault("GCS_BUCKET_LOCATION", "US"),
			StorageClass:     getEnvOrDefault("GCS_STORAGE_CLASS", "COLDLINE"), // Default to COLDLINE

		},
		IPS: IPSConfig{
			Host:     getEnvOrDefault("IPS_HOST", "localhost"),
			Port:     int(getEnvAsInt("IPS_PORT", 8081)),
			Username: getEnvOrDefault("IPS_USERNAME", "admin"),
			Password: getEnvOrDefault("IPS_PASSWORD", "admin"),
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
