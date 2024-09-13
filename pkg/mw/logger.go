package mw

import (
	"net/http"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var RequestLogger zerolog.Logger
var ResponseLogger zerolog.Logger

func init() {
	requestLogFile := &lumberjack.Logger{
		Filename:   "logs/request.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}
	responseLogFile := &lumberjack.Logger{
		Filename:   "logs/response.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}
	RequestLogger = zerolog.New(requestLogFile).With().Timestamp().Logger()
	ResponseLogger = zerolog.New(responseLogFile).With().Timestamp().Logger()
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		RequestLogger.Info().Str("method", r.Method).Str("url", r.URL.String()).Msg("request")
		rw := &responseRecorder{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		ResponseLogger.Info().Str("method", r.Method).Str("url", r.URL.String()).Int("status", rw.statusCode).Int("duration", int(duration.Milliseconds())).Msg("response")
	})
}
