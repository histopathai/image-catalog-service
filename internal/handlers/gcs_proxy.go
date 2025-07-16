package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

type GCSProxyHandler struct {
	GCSClient  *storage.Client
	BucketName string
}

func NewGCSProxyHandler(projectID, bucketName string) (*GCSProxyHandler, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSProxyHandler{
		GCSClient:  client,
		BucketName: bucketName,
	}, nil
}

func (h *GCSProxyHandler) ProxyObject(c *gin.Context) {
	objectPath := strings.TrimPrefix(c.Param("objectPath"), "/") // 🔥 düzeltme

	ctx := context.Background()
	rc, err := h.GCSClient.Bucket(h.BucketName).Object(objectPath).NewReader(ctx)
	if err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("object not found: %s", err.Error()))
		return
	}
	defer rc.Close()

	c.Header("Content-Type", rc.ContentType())
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, rc)
}
