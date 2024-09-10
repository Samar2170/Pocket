package handlers

import (
	"bytes"
	"net/http"
	"pocket/internal"
	"pocket/internal/models"
	"pocket/pkg/response"

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
		response.InternalServerErrorResponse(w, err.Error())
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
