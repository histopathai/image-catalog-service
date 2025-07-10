package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server     ServerConfig
	Firebase   FirebaseConfig
	GCS        GCSConfig
	Processing ProcessingConfig
}

type GCSConfig struct {
	ProjectID         string
	OriginalBucket    string //Cold storage bucket for original images
	DZIBucket         string //Hot storage bucket for DZI files
	ServiceAccountKey string //Path to the service account key file
}

type ProcessingConfig struct {
	MaxFileSize      int64
	SupportedFormats []string
	ThumbnailSize    int
	DZITileSize      int
}

type FirebaseConfig struct {
	ProjectID         string
	ServiceAccountKey string
}

type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() (*Config, error) {
	readTimeout, _ := time.ParseDuration(getEnvOrDefault("READ_TIMEOUT", "30s"))
	writeTimeout, _ := time.ParseDuration(getEnvOrDefault("WRITE_TIMEOUT", "30s"))

	maxFileSize, _ := strconv.ParseInt(getEnvOrDefault("MAX_FILE_SIZE", "104857600"), 10, 64) // 100MB
	thumbnailSize, _ := strconv.Atoi(getEnvOrDefault("THUMBNAIL_SIZE", "256"))
	dziTileSize, _ := strconv.Atoi(getEnvOrDefault("DZI_TILE_SIZE", "256"))

	supportedFormats := strings.Split(getEnvOrDefault("SUPPORTED_FORMATS", "jpg,jpeg,png,tiff,tif,svs,ndpi,mrxs"), ",")

	return &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8080"),
			Environment:  getEnvOrDefault("ENV", "development"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		Firebase: FirebaseConfig{
			ProjectID:         os.Getenv("FIREBASE_PROJECT_ID"),
			ServiceAccountKey: os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"),
		},
		GCS: GCSConfig{
			ProjectID:         os.Getenv("GCS_PROJECT_ID"),
			OriginalBucket:    os.Getenv("GCS_ORIGINAL_BUCKET"),
			ServiceAccountKey: os.Getenv("GCS_SERVICE_ACCOUNT_KEY"),
		},
		Processing: ProcessingConfig{
			MaxFileSize:      maxFileSize,
			SupportedFormats: supportedFormats,
			ThumbnailSize:    thumbnailSize,
			DZITileSize:      dziTileSize,
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) IsSupportedFormat(ext string) bool {
	for _, supported := range c.Processing.SupportedFormats {
		if ext == supported {
			return true
		}
	}
	return false
}
