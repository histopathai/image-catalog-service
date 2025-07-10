package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/histopathai/image-catalog-service/internal/models"
	"github.com/histopathai/image-catalog-service/internal/service"
)

type ImageHandler struct {
	imageService *service.ImageService
}

func NewImageHandler(imgService *service.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imgService,
	}
}

// UploadImage Hadnlers image uploads
func (h *ImageHandler) UploadImage(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")     // Get user ID from Auth-Service
	userRole := c.GetHeader("X-User-Role") // Get user role from Auth-Service

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id_missing", "message": "User ID not found in request headers."})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": "You do not have permission to perform this action."})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_upload_error", "message": "Failed to upload file."})
		return
	}

	var diseaseType *string
	if dt := c.PostForm("disease_type"); dt != "" {
		diseaseType = &dt
	}
	var classification *string
	if cl := c.PostForm("classification"); cl != "" {
		classification = &cl
	}
	var subtype *string
	if st := c.PostForm("subtype"); st != "" {
		subtype = &st
	}

	var grade *string
	if gr := c.PostForm("grade"); gr != "" {
		grade = &gr
	}

	request := &models.ImageUploadRequest{
		DatasetName:    c.PostForm("dataset_name"),
		OrganType:      c.PostForm("organ_type"),
		OriginalUID:    c.PostForm("original_uid"),
		DiseaseType:    diseaseType,
		Classification: classification,
		Subtype:        subtype,
		Grade:          grade,
	}

	image, err := h.imageService.UploadImage(c.Request.Context(), file, header, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image_upload_error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "image": image})
}

// GetImageByID retrieves an image by its ID.
func (h *ImageHandler) GetImageByID(c *gin.Context) {
	imageId := c.Param("image_id")
	if imageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_id_missing", "message": "Image ID is required."})
		return
	}
	image, err := h.imageService.GetImage(c.Request.Context(), imageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image_retrieval_error", "message": err.Error()})
		return
	}
	if image == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image_not_found", "message": "Image not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"image": image})
}

// UpdateImageByID updates an existing image record.
func (h *ImageHandler) UpdateImageByID(c *gin.Context) {
	imageId := c.Param("image_id")
	if imageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_id_missing", "message": "Image ID is required."})
		return
	}

	var updateRequest models.ImageUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Invalid request body."})
		return
	}

	image, err := h.imageService.UpdateImage(c.Request.Context(), imageId, &updateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image_update_error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully", "image": image})
}

// DeleteImageByID deletes an image record and its associated files.
func (h *ImageHandler) DeleteImageByID(c *gin.Context) {
	userID := c.GetHeader("X-User-ID") // Get user ID from Auth-Service
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id_missing", "message": "User ID not found in request headers."})
		return
	}

	role := c.GetHeader("X-User-Role") // Get user role from Auth-Service
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "message": "You do not have permission to perform this action."})
		return
	}

	imageId := c.Param("image_id")
	if imageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image_id_missing", "message": "Image ID is required."})
		return
	}
	err := h.imageService.DeleteImage(c.Request.Context(), imageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image_deletion_error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// GetImages retrieves a list of images with optional filtering.
func (h *ImageHandler) GetImages(c *gin.Context) {
	datasetName := c.Query("dataset_name")
	organType := c.Query("organ_type")
	diseaseType := c.Query("disease_type")
	classification := c.Query("classification")
	subtype := c.Query("subtype")
	grade := c.Query("grade")

	filter := &models.ImageFilter{
		DatasetName:    datasetName,
		OrganType:      organType,
		DiseaseType:    diseaseType,
		Classification: classification,
		Subtype:        subtype,
		Grade:          grade,
	}

	images, err := h.imageService.ListImages(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image_retrieval_error", "message": err.Error()})
		return
	}
	if len(images) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No images found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"images": images})
}
