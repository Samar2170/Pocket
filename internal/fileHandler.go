package internal

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"pocket/internal/models"
	"pocket/pkg/db"
	"pocket/pkg/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	randomString := utils.GenerateRandomString(10)
	fileNameWoExt := strings.Split(fileHeader.Filename, ".")
	fileName, extension := fileNameWoExt[0], fileNameWoExt[1]
	newFileName := fileName + "_" + randomString + "." + extension
	if _, ok := ValidExtensions[extension]; !ok {
		return "", errors.New("unallowed file extension")
	}
	newFilePath := filepath.Join(UploadDir, "pocketstorage", newFileName)
	newFile, err := os.Create(newFilePath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	reader := bufio.NewReader(file)
	writer := io.Writer(newFile)
	_, err = reader.WriteTo(writer)
	if err != nil {
		return "", err
	}
	fileSize := reader.Size()
	fileSizeInMB := utils.ConvertFileSize(float64(fileSize), "bytes", "mb")

	fmd := models.FileMetaData{
		ID:          uuid.New().String(),
		OgFileName:  fileName,
		NewFileName: newFileName,
		FilePath:    newFilePath,
		Extension:   extension,
		Size:        fileSize,
		SizeInMB:    fileSizeInMB,
		CreatedAt:   time.Now(),
	}
	err = db.DB.Create(&fmd).Error
	if err != nil {
		return "", err
	}
	return fmd.ID, nil
}

func SaveImageTelegram(fileURL string) (string, error) {
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

func GetFileByID(id string) ([]byte, models.FileMetaData, error) {
	var fmd models.FileMetaData
	var data []byte
	var err error

	err = db.DB.Where("id = ?", id).First(&fmd).Error
	if err != nil {
		return nil, fmd, err
	}
	data, err = os.ReadFile(fmd.FilePath)
	if err != nil {
		return nil, fmd, err
	}
	return data, fmd, nil
}
