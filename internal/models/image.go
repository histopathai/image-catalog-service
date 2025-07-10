package models

import (
	"time"
)

type Image struct {
	ID             string  `json:"id" firestore:"id"`
	OriginalName   string  `json:"original_name" firestore:"original_name"`
	OriginalUID    string  `json:"original_uid" firestore:"original_uid"`
	DatasetName    string  `json:"dataset_name" firestore:"dataset_name"`
	ImageType      string  `json:"image_type" firestore:"image_type"`
	OrganType      string  `json:"organ_type" firestore:"organ_type"`
	DiseaseType    *string `json:"disease_type,omitempty" firestore:"disease_type,omitempty"`
	Classification *string `json:"classification,omitempty" firestore:"classification,omitempty"`
	SubType        *string `json:"sub_type,omitempty" firestore:"sub_type,omitempty"`
	Grade          *string `json:"grade,omitempty" firestore:"grade,omitempty"`

	// File paths
	OriginalGCSPath  string `json:"original_gcs_path" firestore:"original_gcs_path"`
	DZIGCSPath       string `json:"dzi_gcs_path" firestore:"dzi_gcs_path"`
	TilesGCSPath     string `json:"tiles_gcs_path" firestore:"tiles_gcs_path"`
	ThumbnailGCSPath string `json:"thumbnail_gcs_path" firestore:"thumbnail_gcs_path"`

	// Metadata
	Width            int    `json:"width" firestore:"width"`
	Height           int    `json:"height" firestore:"height"`
	FileSize         int64  `json:"file_size" firestore:"file_size"`
	ProcessingStatus string `json:"processing_status" firestore:"processing_status"` // pending, processing, completed, failed

	// Timestamps
	CreatedAt   time.Time  `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" firestore:"updated_at"`
	ProcessedAt *time.Time `json:"processed_at,omitempty" firestore:"processed_at,omitempty"`
}

type ImageUploadRequest struct {
	DatasetName    string  `json:"dataset_name" validate:"required"`
	OrganType      string  `json:"organ_type" validate:"required"`
	OriginalUID    string  `json:"original_uid" validate:"required"`
	DiseaseType    *string `json:"disease_type,omitempty"`
	Classification *string `json:"classification,omitempty"`
	Subtype        *string `json:"subtype,omitempty"`
	Grade          *string `json:"grade,omitempty"`
}

type ImageUpdateRequest struct {
	DatasetName    *string `json:"dataset_name,omitempty"`
	OrganType      *string `json:"organ_type,omitempty"`
	DiseaseType    *string `json:"disease_type,omitempty"`
	Classification *string `json:"classification,omitempty"`
	OriginalUID    *string `json:"original_uid,omitempty"`
	Subtype        *string `json:"subtype,omitempty"`
	Grade          *string `json:"grade,omitempty"`
}

type ImageFilter struct {
	DatasetName    string `json:"dataset_name,omitempty"`
	OrganType      string `json:"organ_type,omitempty"`
	DiseaseType    string `json:"disease_type,omitempty"`
	Classification string `json:"classification,omitempty"`
	Subtype        string `json:"subtype,omitempty"`
	Grade          string `json:"grade,omitempty"`
}
