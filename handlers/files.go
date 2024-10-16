package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"pocket/internal"
	"pocket/internal/models"
	"pocket/pkg/response"
	"time"

	"github.com/gorilla/mux"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.BadRequestResponse(w, err.Error())
		return
	}
	url, err := internal.SaveFile(file, fileHeader)
	if err != nil {
		response.BadRequestResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]string{"file-id": url})
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["id"]
	data, fmd, err := internal.GetFileByID(fileId)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	http.ServeContent(w, r, fmd.FilePath, fmd.CreatedAt, bytes.NewReader(data))
}

type FileMetaDataResponse struct {
	Name      string
	ID        string
	Extension string
	SizeInMB  float64
}

func GetFileMetaDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["id"]
	fmd, err := models.GetFileMetaDataById(fileId)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	fmdr := FileMetaDataResponse{
		Name:      fmd.OgFileName,
		ID:        fmd.ID,
		Extension: fmd.Extension,
		SizeInMB:  fmd.SizeInMB,
	}
	response.JSONResponse(w, fmdr)
}

type SyncFilesRequest struct {
	CutoffTime string `json:"cutoff_time"`
}

func SyncFilesHandler(w http.ResponseWriter, r *http.Request) {
	//project is master server, only serve files syncing

	var sfr SyncFilesRequest
	err := json.NewDecoder(r.Body).Decode(&sfr)
	if err != nil {
		response.BadRequestResponse(w, "Error parsing request  "+err.Error())
		return
	}
	cutoffTimeParsed, err := time.Parse("2006-01-02T15:04:05Z", sfr.CutoffTime)
	if err != nil {
		response.BadRequestResponse(w, "Error parsing time  "+err.Error())
		return
	}
	fileName, err := internal.DownloadFiles(cutoffTimeParsed)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error downloading files  "+err.Error())
		return
	}
	file, err := os.Open(internal.TMPFOLDER + "/" + fileName)
	if err != nil {
		response.InternalServerErrorResponse(w, "Error opening file  "+err.Error())
		return
	}
	defer file.Close()
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	_, err = io.Copy(w, file)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}

	os.Remove(internal.TMPFOLDER + "/" + fileName)
}
