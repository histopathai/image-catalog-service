package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/histopathai/image-catalog-service/internal/models"
)

func FileInfo(req *models.JobRequest) (*models.Image, error) {

	filepath := req.ImagePath
	// exist, err :=
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filepath)
	}

	ext := req.Image.ImageType

	if ext == ".svs" {
		return getSVSInfo(req)
	} else {
		return getVIPSInfo(req)
	}

}

func ExportThumbnail(filepath, ext, outputPath string, thumbSize int) error {
	if ext == ".svs" {
		return exportSVSThumbnail(filepath, outputPath, thumbSize)
	} else {
		return exportVIPSThumbnail(filepath, outputPath, thumbSize)
	}
}

func getSVSInfo(req *models.JobRequest) (*models.Image, error) {
	cmd := exec.Command("openslide-show-properties", req.ImagePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute openslide command: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.SplitN(strings.TrimSpace(line), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		switch key {
		case "openslide.level[0].width":
			if _, err := fmt.Sscanf(value, "%d", &req.Image.Width); err != nil {
				return nil, fmt.Errorf("invalid width value: %w", err)
			}
		case "openslide.level[0].height":
			if _, err := fmt.Sscanf(value, "%d", &req.Image.Height); err != nil {
				return nil, fmt.Errorf("invalid height value: %w", err)
			}
		}
	}
	return &req.Image, nil
}

func getVIPSInfo(req *models.JobRequest) (*models.Image, error) {
	cmd := exec.Command("vipsheader", "-f", "json", req.ImagePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute vipsheader: %w", err)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(output, &header); err != nil {
		return nil, fmt.Errorf("failed to parse vipsheader output: %w", err)
	}

	if width, ok := header["Xsize"].(float64); ok {
		req.Image.Width = int(width)
	}
	if height, ok := header["Ysize"].(float64); ok {
		req.Image.Height = int(height)
	}
	return &req.Image, nil
}

func exportSVSThumbnail(inputPath, outputPath string, thumbSize int) error {
	// 1. Geçici PNG dosyası
	tempPng, err := os.CreateTemp("", "thumb-*.png")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempPng.Name()) // iş bitince sil

	level, err := getDeepestSVSLevel(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get deepest SVS level: %w", err)
	}

	cmd := exec.Command("openslide-write-png", inputPath, tempPng.Name(),
		"--level", level)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openslide-write-png failed: %w", err)
	}

	cmd = exec.Command("vips", "thumbnail", tempPng.Name(), outputPath,
		strconv.Itoa(thumbSize),
		"--size", "both", "--format", "jpg")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("vips thumbnail failed: %w", err)
	}

	return nil
}

func exportVIPSThumbnail(inputPath, outputPath string, thumbSize int) error {
	cmd := exec.Command("vips", "thumbnail", inputPath, outputPath,
		strconv.Itoa(thumbSize),
		"--size", "both")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create thumbnail with vips: %w", err)
	}

	return nil
}

func getDeepestSVSLevel(filepath string) (string, error) {
	cmd := exec.Command("openslide-show-properties", filepath)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	maxLevel := -1
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "openslide.level[") && strings.Contains(line, "width") {
			var level int
			if _, err := fmt.Sscanf(line, "openslide.level[%d].width", &level); err == nil {
				if level > maxLevel {
					maxLevel = level
				}
			}
		}
	}

	if maxLevel == -1 {
		return "", fmt.Errorf("no level found")
	}
	return strconv.Itoa(maxLevel), nil
}

func ExtractDZI(inputpath, outputPath string, request *models.JobRequest) error {

	params := request.Parameters
	tileSize := params.TileSize
	overlap := params.Overlap
	quality := params.Quality
	suffix := params.Suffix
	layout := params.Layout

	if tileSize <= 0 {
		return errors.New("tile_size must be a positive integer")
	}

	if overlap < 0 {
		return errors.New("overlap must be a non-negative integer")
	}

	if overlap >= tileSize {
		return errors.New("overlap must be less than tile_size")
	}

	if suffix == ".jpg" || suffix == ".jpeg" {
		if quality < 0 || quality > 100 {
			return errors.New("quality for JPEG must be between 0 and 100")
		}
	} else if quality != 0 {
		fmt.Printf("Warning: Quality parameter is ignored for non-JPEG formats, using default value.\n")
	}

	supportedSuffixes := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}

	if !supportedSuffixes[strings.ToLower(suffix)] {
		return fmt.Errorf("unsupported suffix: %s. Supported formats are .jpg, .jpeg, .png", suffix)
	}

	switch strings.ToLower(layout) {
	case "dzi":
		// DZI layout, no action needed
	case "google":
		if suffix == "png" {
			fmt.Printf("Warning: Google layout does not support PNG suffix, using DZI layout instead.\n")
			layout = "dzi"
		}
	default:
		return fmt.Errorf("unsupported layout: %s. Supported layouts are 'dzi' and 'google'", layout)
	}

	args := []string{
		"dzsave",
		inputpath,
		outputPath,
		"--layout", layout,
		"--tile-size", fmt.Sprintf("%d", tileSize),
		"--overlap", fmt.Sprintf("%d", overlap),
		"--suffix", suffix,
	}

	if suffix == ".jpg" || suffix == ".jpeg" {
		args = append(args, "--quality", fmt.Sprintf("%d", quality))
	}

	cmd := exec.Command("vips", args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to create DZI: %w - VIPS Output: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
