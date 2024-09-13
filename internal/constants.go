package internal

import (
	"os"

	"pocket/pkg/auditlog"

	"github.com/joho/godotenv"
)

var UploadDir string

const (
	SUBFOLDER = "pocketstorage"
	TMPFOLDER = "tmp"
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
		auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error loading .env file")
	}
	UploadDir = os.Getenv("UPLOADDIR")
}
