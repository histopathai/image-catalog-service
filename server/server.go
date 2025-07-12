package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/histopathai/image-catalog-service/config"
	"github.com/histopathai/image-catalog-service/internal/handlers"
	"github.com/histopathai/image-catalog-service/internal/routes"
)

type Server struct {
	httpServer *http.Server
	config     *config.Config
}

func NewServer(cfg *config.Config, imageHandler *handlers.ImageHandler) *Server {

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := routes.SetupRouter(imageHandler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &Server{
		httpServer: httpServer,
		config:     cfg,
	}
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.config.Server.Port)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//Shut down the server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		return err
	}

	slog.Info("Server gracefully stopped")
	return nil
}
