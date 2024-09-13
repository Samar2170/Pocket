package auditlog

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var Errorlogger zerolog.Logger

func init() {
	logFile := &lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}
	Errorlogger = zerolog.New(logFile).With().Timestamp().Logger()
}
