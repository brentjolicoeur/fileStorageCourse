package main

import (
	"errors"
	"mime"
	"net/http"
	"strings"
)

func determineFileExtension(header http.Header) (string, error) {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		return contentType, errors.New("Missing Content-Type Header")
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", errors.New("Error parsing media type")
	}
	_, fileExtension, ok := strings.Cut(mediaType, "/")
	if !ok {
		return "", errors.New("error determining file extension")
	}
	if fileExtension == "jpeg" {
		fileExtension = "jpg"
	}

	return fileExtension, nil
}
