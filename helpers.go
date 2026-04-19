package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"mime"
	"strings"
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
