package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"github.com/histopathai/image-catalog-service/adapter"
	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/handlers"
	"github.com/histopathai/image-catalog-service/internal/service"
	"github.com/histopathai/image-catalog-service/server"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	// Initialize context
	ctx := context.Background()

	fmt.Printf("Loaded configuration: %+v\n", cfg)

	// Initialize Firestore
	firestoreClient, err := initFireStore(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize Firestore", "error", err)
		os.Exit(1)
	}

	// Initialize ImageService
	imageService, err := initImageService(firestoreClient, cfg)

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

	gcsProxyHandler, err := handlers.NewGCSProxyHandler(cfg.ProjectID, cfg.BucketName)
	if err != nil {
		slog.Error("Failed to create GCSProxyHandler", "error", err)
		os.Exit(1)
	}

	// Initialize Server
	server := server.NewServer(cfg, imageHandler, gcsProxyHandler)

	if server == nil {
		slog.Error("Failed to create Server")
		os.Exit(1)
	}

	// Start the server
	if err := server.Start(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	slog.Info("Image processing result subscriber started")
}

func initImageService(firestoreClient *firestore.Client, cfg *config.Config) (*service.ImageService, error) {
	if firestoreClient == nil {
		return nil, fmt.Errorf("firestore client is nil")
	}

	repo, err := adapter.NewFirestoreCollection(firestoreClient, "images")
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore repository: %w", err)
	}

	imageService := service.NewImageService(repo, cfg)
	if imageService == nil {
		return nil, fmt.Errorf("failed to create ImageService")
	}

	return imageService, nil
}

func initFireStore(ctx context.Context, cfg *config.Config) (*firestore.Client, error) {
	var app *firebase.App
	var err error

	app, err = firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.ProjectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Firebase app: %w", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return firestoreClient, nil
}
