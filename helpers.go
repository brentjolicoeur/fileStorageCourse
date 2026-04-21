package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"mime"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func determineFileExtension(contentType string) (string, error) {

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("Error parsing media type: %w", err)
	}
	_, fileExtension, ok := strings.Cut(mediaType, "/")
	if !ok {
		return "", fmt.Errorf("error determining file extension: %w", err)
	}
	if fileExtension == "jpeg" {
		fileExtension = "jpg"
	}

	return fileExtension, nil
}

func randomKeyName() (string, error) {
	randomName := make([]byte, 32)
	_, err := rand.Read(randomName)
	if err != nil {
		return "", fmt.Errorf("error generating random key name: %w", err)
	}
	rawString := base64.RawURLEncoding.EncodeToString(randomName)

	return rawString, nil
}

func (cfg apiConfig) getObjectURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
}

func getVideoAspectRatio(filePath string) (string, error) {
	var data bytes.Buffer

	command := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	command.Stdout = &data

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("Error running ffprobe command: %w", err)
	}

	type Stream struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type FFProbeOutput struct {
		Streams []Stream `json:"streams"`
	}

	var ffprobeOutput FFProbeOutput

	err = json.Unmarshal(data.Bytes(), &ffprobeOutput)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	if len(ffprobeOutput.Streams) == 0 {
		return "", fmt.Errorf("no streams found in file: %s", filePath)
	}

	width := ffprobeOutput.Streams[0].Width
	height := ffprobeOutput.Streams[0].Height

	ratio := float64(width) / float64(height)

	if math.Abs(ratio-9.0/16.0) < 0.01 {
		return "portrait", nil
	} else if math.Abs(ratio-16.0/9.0) < 0.01 {
		return "landscape", nil
	} else {
		return "other", nil
	}
}

func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3Client)

	params := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	req, err := presignClient.PresignGetObject(context.Background(), &params, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", fmt.Errorf("Error Getting Presigned Request: %w", err)
	}

	return req.URL, nil
}
