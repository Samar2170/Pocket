package internal

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func SaveImage(fileURL string) (string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fileName := fileURL[strings.LastIndex(fileURL, "/")+1:]
	err = os.WriteFile("uploads/"+fileName, byteValue, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func SaveImageThumbnail(fileURL string) (string, error) {
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fileName := "thumbnail-" + fileURL[strings.LastIndex(fileURL, "/")+1:]
	err = os.WriteFile("uploads/"+fileName, byteValue, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
