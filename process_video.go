package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	outputPath := filePath + ".processing"

	command := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputPath)

	err := command.Run()
	if err != nil {
		return "", fmt.Errorf("Error running ffmpeg command: %w", err)
	}

	return outputPath, nil
}
