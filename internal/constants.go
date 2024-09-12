package internal

import (
	"os"

	"github.com/joho/godotenv"
)

var UploadDir string

const (
	SUBFOLDER = "pocketstorage"
)

var ValidExtensions = map[string]struct{}{
	"pdf":  {},
	"docx": {},
	"doc":  {},

	"jpg":  {},
	"jpeg": {},
	"png":  {},

	"mp4": {},

	"xlsx": {},
	"xls":  {},
	"csv":  {},
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	UploadDir = os.Getenv("UPLOADDIR")
}
