package main

import (
	"net/http"
	"pocket/handlers"

	"github.com/gorilla/mux"
)

func main() {

	// args := os.Args[1:]
	// switch args[0] {
	// case "fxb":
	// 	RunFXBStorageServer()
	// case "storage":
	// 	RunStorageServer()
	// case "server":
	// 	RunTelegramServer()
	// default:
	RunTelegramServer()

	// }
}

func RunStorageServer() {
	mux := mux.NewRouter()
	storage := mux.PathPrefix("/storage").Subrouter()
	storage.Handle("/health/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	getFileHandler := http.HandlerFunc(handlers.GetFileHandler)
	storage.Handle("/file/{id}", getFileHandler).Methods("GET")
	getFileMetaDataHandler := http.HandlerFunc(handlers.GetFileMetaDataHandler)
	storage.Handle("/metadata/{id}", getFileMetaDataHandler).Methods("GET")

	uploadFileHandler := http.HandlerFunc(handlers.UploadFileHandler)
	storage.Handle("/upload", uploadFileHandler).Methods("POST")
	http.ListenAndServe(":8080", mux)
}
