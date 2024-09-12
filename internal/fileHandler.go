package internal

import (
	"archive/zip"
	"bufio"
	"bytes"
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
	"gorm.io/gorm"
)

func SaveFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	randomString := utils.GenerateRandomString(10)
	fileNameWoExt := strings.Split(fileHeader.Filename, ".")
	fileName, extension := fileNameWoExt[0], fileNameWoExt[1]
	newFileName := fileName + "_" + randomString + "." + extension
	if _, ok := ValidExtensions[extension]; !ok {
		return "", errors.New("unallowed file extension")
	}
	newFilePath := filepath.Join(UploadDir, SUBFOLDER, newFileName)
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
	newFilePath := filepath.Join(UploadDir, SUBFOLDER, fileName)
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
		var existingTag models.Tag
		var newTag models.Tag
		err := db.DB.Where("name = ?", tag).First(&existingTag).Error
		if err == gorm.ErrRecordNotFound {
			newTag = models.Tag{
				Name: tag,
			}
			err = db.DB.Create(&newTag).Error
			if err != nil {
				return err
			}
			existingTag = newTag
		}
		if err != nil {
			return err
		}
		ft := models.FileTag{
			ID:     uuid.New().String(),
			FileID: fileId,
			TagID:  existingTag.ID,
		}
		err = db.DB.Create(&ft).Error
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

func DownloadFiles(cutoff time.Time) (string, error) {
	var buf bytes.Buffer
	var fmds []models.FileMetaData
	var err error
	var opFileName string
	zipWriter := zip.NewWriter(&buf)

	err = db.DB.Where("created_at > ?", cutoff).Find(&fmds).Error
	if err != nil {
		return "", errors.New("Error fetching files " + err.Error())
	}
	opFileName = "archive_" + utils.GenerateRandomString(5) + ".zip"

	for _, fmd := range fmds {
		content, err := os.ReadFile(fmd.FilePath)
		if err != nil {
			return "", err
		}
		f, err := zipWriter.Create(fmd.NewFileName)
		if err != nil {
			return "", err
		}
		_, err = f.Write(content)
		if err != nil {
			return "", err
		}
	}
	err = zipWriter.Close()
	if err != nil {
		return "", err
	}
	f, err := os.Create(TMPFOLDER + "/" + opFileName)
	if err != nil {
		return "", err
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return "", err
	}
	return opFileName, nil
}
