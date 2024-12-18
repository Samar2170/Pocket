package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"pocket/handlers"
	"pocket/internal"
	"pocket/pkg/auditlog"
	"pocket/pkg/auth"
	"pocket/pkg/mw"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "setup":
			Setup()
		case "fxb":
			RunFXBStorageServer()
		case "storage":
			RunStorageServer()
		case "keygen":
			auth.GetNewKey()
		case "skey":
			key := auth.GenerateKey(32)
			auditlog.AuditLogger.Info().Str("key", key).Msg("Generated key")
		}
	} else {
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
	auth.GetNewKey()
}

func RunStorageServer() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	var networkIp string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ipnet.IP.String()[:4] == "192" {
					networkIp = ipnet.IP.String()
					break
				}
			}
		}
	}
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

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
	auditlog.AuditLogger.Println("Storage server started at " + networkIp + ":" + Port)
	srv := &http.Server{
		Handler:      wrappedMux,
		Addr:         networkIp + ":" + Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
