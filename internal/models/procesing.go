package models

type ProcessingJob struct {
	ImageID           string `json:"image_id" firestore:"image_id"`
	OriginalGCSPath   string `json:"original_gcs_path" firestore:"original_gcs_path"`
	FileType          string `json:"file_type" firestore:"file_type"`
	DestinationBucket string `json:"destination_bucket" firestore:"destination_bucket"`
	ThumbnailSize     int    `json:"thumbnail_size" firestore:"thumbnail_size"`
	DZITileSize       int    `json:"dzi_tile_size" firestore:"dzi_tile_size"`
}

const (
	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusCompleted  = "completed"
	ProcessingStatusFailed     = "failed"
	ProcessingStatusCancelled  = "cancelled"
)
