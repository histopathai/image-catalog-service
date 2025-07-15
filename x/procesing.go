package models

import "time" // Needed for CreatedAt/UpdatedAt in JobParameters

type JobParameters struct {
	ID            string    `json:"id" firestore:"id"`
	ThumbnailSize int       `json:"thumbnail_size" firestore:"thumbnail_size"`
	TileSize      int       `json:"tile_size" firestore:"tile_size"`
	Overlap       int       `json:"overlap" firestore:"overlap"`
	Quality       int       `json:"quality" firestore:"quality"`
	Layout        string    `json:"layout" firestore:"layout"`
	Suffix        string    `json:"suffix" firestore:"suffix"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" firestore:"updated_at"`
	IsDefault     bool      `json:"is_default" firestore:"is_default"`
}

type ProcessingJob struct {
	ImageID       string        `json:"image_id" firestore:"image_id"`
	OriginalPath  string        `json:"original_path" firestore:"original_path"`
	ImageType     string        `json:"image_type" firestore:"image_type"`
	JobParameters JobParameters `json:"job_parameters" firestore:"job_parameters"`
}

type ProcessingResultJob struct {
	ImageID                  string `json:"image_id" firestore:"image_id"`
	Status                   string `json:"status" firestore:"status"`
	ErrorMessage             string `json:"error_message,omitempty" firestore:"error_message,omitempty"`
	GCSThumbnailPath         string `json:"gcs_thumbnail_path,omitempty" firestore:"gcs_thumbnail_path,omitempty"`
	GCSDZIPath               string `json:"gcs_dzi_path,omitempty" firestore:"gcs_dzi_path,omitempty"`
	GCSTilePath              string `json:"gcs_tile_path,omitempty" firestore:"gcs_tile_path,omitempty"`
	Width                    int    `json:"width,omitempty" firestore:"width,omitempty"`
	Height                   int    `json:"height,omitempty" firestore:"height,omitempty"`
	OriginalLocalCleanupPath string `json:"original_local_cleanup_path,omitempty"`
}

const (
	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusCompleted  = "completed"
	ProcessingStatusFailed     = "failed"
)
