package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pocket/handlers"
	"pocket/internal"
	"pocket/pkg/mw"
	"sync"

	"github.com/gorilla/mux"
)

func main() {

	args := os.Args[1:]
	switch args[0] {
	case "setup":
		Setup()
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

func Setup() {
	uploadDir := internal.UploadDir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		err := os.MkdirAll("tmp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.MkdirAll("logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
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
	wrappedMux := mw.LogRequest(mux)
	addr := fmt.Sprintf("%s:%s", Host, Port)
	http.ListenAndServe(addr, wrappedMux)
}
