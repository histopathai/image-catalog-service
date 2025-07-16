package models

import (
	"time"
)

type Image struct {
	ID             string  `json:"id" firestore:"id"`
	FileName       string  `json:"file_name" firestore:"file_name"`
	FileUID        string  `json:"file_uid" firestore:"file_uid"`
	DatasetName    string  `json:"dataset_name" firestore:"dataset_name"`
	OrganType      string  `json:"organ_type" firestore:"organ_type"`
	DiseaseType    *string `json:"disease_type,omitempty" firestore:"disease_type,omitempty"`
	Classification *string `json:"classification,omitempty" firestore:"classification,omitempty"`
	SubType        *string `json:"sub_type,omitempty" firestore:"sub_type,omitempty"` // Consistent casing
	Grade          *string `json:"grade,omitempty" firestore:"grade,omitempty"`

	DZIGCSPath       string `json:"dzi_gcs_path" firestore:"dzi_gcs_path"`
	TilesGCSPath     string `json:"tiles_gcs_path" firestore:"tiles_gcs_path"`
	ThumbnailGCSPath string `json:"thumbnail_gcs_path" firestore:"thumbnail_gcs_path"`

	// Metadata
	Width  int    `json:"width" firestore:"width"`
	Height int    `json:"height" firestore:"height"`
	Size   int64  `json:"size" firestore:"size"`
	Format string `json:"format"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

type ImageFilter struct {
	DatasetName    *string `json:"dataset_name,omitempty" firestore:"dataset_name,omitempty"`
	OrganType      *string `json:"organ_type,omitempty" firestore:"organ_type,omitempty"`
	DiseaseType    *string `json:"disease_type,omitempty" firestore:"disease_type,omitempty"`
	Classification *string `json:"classification,omitempty" firestore:"classification,omitempty"`
	SubType        *string `json:"sub_type,omitempty" firestore:"sub_type,omitempty"`
	Grade          *string `json:"grade,omitempty" firestore:"grade,omitempty"`
}

type ImageUpdateRequest ImageFilter
