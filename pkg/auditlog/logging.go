package auditlog

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var Errorlogger zerolog.Logger
var AuditLogger zerolog.Logger
var mode string = "dev"

func init() {
	godotenv.Load(".env")
	if os.Getenv("MODE") != "" {
		mode = os.Getenv("MODE")
	}
	if mode == "dev" {
		AuditLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		Errorlogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		auditLogFile := &lumberjack.Logger{
			Filename:   "logs/audit.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   false,
		}
		logFile := &lumberjack.Logger{
			Filename:   "logs/error.log",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   false,
		}
		Errorlogger = zerolog.New(logFile).With().Timestamp().Logger()
		AuditLogger = zerolog.New(auditLogFile).With().Timestamp().Logger()

	}
}
