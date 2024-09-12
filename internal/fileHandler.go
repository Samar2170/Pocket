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
	"strconv"
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

func SaveFileTelegram(fileURL string) (string, error) {
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
	newFilePath := filepath.Join(UploadDir, "pocketstorage", fileName)
	err = os.WriteFile(newFilePath, byteValue, 0644)
	if err != nil {
		return "", err
	}
	newFileName := fileName
	var conflicts int64
	db.DB.Model(&models.FileMetaData{}).Where("new_file_name = ?", fileName).Count(&conflicts)
	if conflicts > 0 {
		tmpF := strings.Split(fileName, ".")
		tmpF[0] = tmpF[0] + "_" + strconv.Itoa(int(conflicts))
		newFileName = tmpF[0] + "." + tmpF[len(tmpF)-1]
	}
	fileId := uuid.New().String()
	fmd := models.FileMetaData{
		ID:          fileId,
		OgFileName:  fileName,
		NewFileName: newFileName,
		FilePath:    newFilePath,
		Extension:   filepath.Ext(fileName),
		Size:        len(byteValue),
		SizeInMB:    utils.ConvertFileSize(float64(len(byteValue)), "bytes", "mb"),
		CreatedAt:   time.Now(),
	}
	err = db.DB.Create(&fmd).Error
	if err != nil {
		return "", err
	}
	return fileId, nil
}

func SaveFileCaption(fileId, caption string) error {
	fc := models.FileCaption{
		ID:      uuid.New().String(),
		FileID:  fileId,
		Caption: caption,
	}
	return db.DB.Create(&fc).Error
}

func SaveFileTags(fileId, tags string) error {
	tagsSplit := strings.Split(tags, " ")
	for _, tag := range tagsSplit {
		ft := models.FileTag{
			ID:     uuid.New().String(),
			FileID: fileId,
			Tag:    tag,
		}
		err := db.DB.Create(&ft).Error
		if err != nil {
			return err
		}
	}
	return nil
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
