package main

import (
	"log"
	"net/http"
	"os"
	"pocket/handlers"
	"sync"

	"github.com/gorilla/mux"
)

func main() {

	args := os.Args[1:]
	switch args[0] {
	case "fxb":
		RunFXBStorageServer()
	case "storage":
		RunStorageServer()
	case "server":
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			RunTelegramServer()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			RunStorageServer()
			wg.Done()
		}()
		wg.Wait()
	default:
		log.Panic("Invalid argument")
	}
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

	syncFilesHandler := http.HandlerFunc(handlers.SyncFilesHandler)
	storage.Handle("/sync", syncFilesHandler).Methods("POST")

	uploadFileHandler := http.HandlerFunc(handlers.UploadFileHandler)
	storage.Handle("/upload", uploadFileHandler).Methods("POST")
	http.ListenAndServe(":8080", mux)
}
