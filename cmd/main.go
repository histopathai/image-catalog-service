package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"

	"github.com/histopathai/image-catalog-service/adapter"
	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/handlers"
	"github.com/histopathai/image-catalog-service/internal/service"
	"github.com/histopathai/image-catalog-service/server"
	"github.com/streadway/amqp"
)

func main() {
	//
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	// Initialize context
	ctx := context.Background()

	// Load configuration
	cfg := loadConfig()

	// Initialize Firestore
	firestoreClient, err := initFireStore(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize Firestore", "error", err)
		os.Exit(1)
	}

	// Initialize GCS
	gcsClient, err := initGCS(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize GCS", "error", err)
		os.Exit(1)
	}

	// Initialize RabbitMQ
	mq, err := initMQ(ctx, cfg)
	if err != nil {
		slog.Error("Failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}

	// Initialize ImageService
	imageService, err := initImageService(ctx, cfg, firestoreClient, gcsClient, mq)
	if err != nil {
		slog.Error("Failed to initialize ImageService", "error", err)
		os.Exit(1)
	}

	// Initialize Handlers
	imageHandler := handlers.NewImageHandler(imageService)
	if imageHandler == nil {
		slog.Error("Failed to create ImageHandler")
		os.Exit(1)
	}

	// Initialize Server
	server := server.NewServer(cfg, imageHandler)

	if server == nil {
		slog.Error("Failed to create Server")
		os.Exit(1)
	}

	// Start the server
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	go imageService.StartProcessingResultSubscriber(ctx)
	slog.Info("Image processing result subscriber started")
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	return cfg
}

func initFireStore(ctx context.Context, cfg *config.Config) (*firestore.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return firestoreClient, nil
}

func initGCS(ctx context.Context, cfg *config.Config) (*storage.Client, error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	bucketName := cfg.GCS.OriginalBucket
	projectID := cfg.Server.PROJECT_ID
	bucket := client.Bucket(bucketName)

	//Check if bucket exists
	_, err = bucket.Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{
			Location:     cfg.GCS.BucketLocation,
			StorageClass: cfg.GCS.StorageClass,
		}); err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to access bucket %s: %w", bucketName, err)
	}

	return client, nil
}

func initMQ(ctx context.Context, cfg *config.Config) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.IPS.Username, cfg.IPS.Password, cfg.IPS.Host, cfg.IPS.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	slog.Info("Connected to RabbitMQ", "host", cfg.IPS.Host, "port", cfg.IPS.Port)
	return conn, nil
}

func initImageService(ctx context.Context, cfg *config.Config, firestoreClient *firestore.Client, gcsClient *storage.Client, mq *amqp.Connection) (*service.ImageService, error) {

	// init
	var jonQueue = "image_processing_jobs"
	var resultQueue = "image_processing_results"
	mq_client, err := service.NewRabbitMQ(mq, jonQueue, resultQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RabbitMQ client: %w", err)
	}

	// Initialize repositories

	imageRepo, err := adapter.NewFirestoreCollection(firestoreClient, "images")
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore image repository: %w", err)
	}

	fileRepo := adapter.NewGCSFileRepository(gcsClient, cfg.GCS.OriginalBucket)

	// Initialize the ImageService

	imageService := service.NewImageService(imageRepo, fileRepo, mq_client, cfg)
	if imageService == nil {
		return nil, fmt.Errorf("failed to create ImageService")
	}
	slog.Info("ImageService initialized successfully")
	return imageService, nil
}
