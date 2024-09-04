package handlers

import (
	"bytes"
	"net/http"
	"pocket/internal"
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
	response.SuccessResponse(w, url)
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileId := vars["fileId"]
	data, fmd, err := internal.GetFileByID(fileId)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	http.ServeContent(w, r, fmd.FilePath, fmd.CreatedAt, bytes.NewReader(data))
}
