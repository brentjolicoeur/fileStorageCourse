package main

import (
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
