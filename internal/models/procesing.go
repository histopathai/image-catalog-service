package models

type ProcessingJob struct {
	ImageID         string `json:"image_id" firestore:"image_id"`
	OriginalGCSPath string `json:"original_gcs_path" firestore:"original_gcs_path"`
	FileType        string `json:"file_type" firestore:"file_type"`
}

type ProcessingResultJob struct {
	ImageID          string `json:"image_id"`
	Status           string `json:"status"`
	GCSThumbnailPath string `json:"gcs_thumbnail_path,omitempty"`
	GCSDZIPath       string `json:"gcs_dzi_path,omitempty"`
	GCSTilePath      string `json:"gcs_file_path,omitempty"`
	Width            int    `json:"width,omitempty"`
	Height           int    `json:"height,omitempty"`
	TileSize         int    `json:"tile_size,omitempty"`
	ThumbnailSize    int    `json:"thumbnail_size,omitempty"`
	Overlap          int    `json:"overlap,omitempty"`
	ErrorMessage     string `json:"error_message,omitempty"`
}

const (
	ProcessingStatusPending    = "pending"
	ProcessingStatusProcessing = "processing"
	ProcessingStatusCompleted  = "completed"
	ProcessingStatusFailed     = "failed"
)
