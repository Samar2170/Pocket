package main

import (
	"log"
	"net/http"
	"os"
	"pocket/handlers"
	"pocket/internal"
	"sync"

	"github.com/gorilla/mux"
)

func main() {

	args := os.Args[1:]
	switch args[0] {
	case "load":
		log.Println("loading accounts")
		err := internal.LoadAccounts()
		if err != nil {
			log.Println(err)
		}
	case "fxb":
		RunFXBStorageServer()
	case "server":
		var wg sync.WaitGroup
		// wg.Add(1)
		// go func() {
		// 	RunFXBStorageServer()
		// 	wg.Done()
		// }()
		wg.Add(1)
		go func() {
			RunTelegramServer()
			wg.Done()
		}()
		wg.Wait()

	default:
		log.Fatalln("invalid argument")
	}
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	RunStorageServer()
	// 	wg.Done()
	// }()

	// wg.Add(1)
	// go func() {
	// 	wg.Done()
	// }()
	// wg.Wait()
}

func RunStorageServer() {
	mux := mux.NewRouter()
	storage := mux.PathPrefix("storage").Subrouter()
	storage.Handle("/health/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	getFileHandler := http.HandlerFunc(handlers.GetFileHandler)
	storage.Handle("/file/{id}", getFileHandler).Methods("GET")

	uploadFileHandler := http.HandlerFunc(handlers.UploadFileHandler)
	storage.Handle("/upload", uploadFileHandler).Methods("POST")
	http.ListenAndServe(":8080", mux)
}
