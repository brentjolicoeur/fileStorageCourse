package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"

	command := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("error running ffmpeg command: %w", err)
	}

	return outputPath, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, nil
	}

	bucket, key, ok := strings.Cut(*video.VideoURL, ",")
	if !ok {
		return database.Video{}, fmt.Errorf("invalid URL format: need <bucket>,<key>")
	}
	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Minute*15)
	if err != nil {
		return database.Video{}, fmt.Errorf("error generating presignedURL: %w", err)
	}

	video.VideoURL = &presignedURL

	return video, nil
}
