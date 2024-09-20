package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pocket/handlers"
	"pocket/internal"
	"pocket/pkg/auditlog"
	"pocket/pkg/auth"
	"pocket/pkg/mw"

	"github.com/gorilla/mux"
)

func main() {
	exePath, err := os.Getwd()
	if err != nil {
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error getting executable path")
		panic(err)
	}
	fmt.Println(exePath)

	// args := os.Args[1:]
	// if len(args) > 0 {
	// 	switch args[0] {
	// 	case "setup":
	// 		Setup()
	// 	case "fxb":
	// 		RunFXBStorageServer()
	// 	case "storage":
	// 		RunStorageServer()
	// 	}
	// } else {
	// 	var wg sync.WaitGroup
	// 	wg.Add(1)
	// 	go func() {
	// 		RunTelegramServer()
	// 		wg.Done()
	// 	}()
	// 	wg.Add(1)
	// 	go func() {
	// 		RunStorageServer()
	// 		wg.Done()
	// 	}()
	// 	wg.Wait()
	// }
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
	auth.GetNewKey()
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
	wrappedMux = mw.APIKeyMiddleware(wrappedMux)
	addr := fmt.Sprintf("%s:%s", Host, Port)
	http.ListenAndServe(addr, wrappedMux)
}
