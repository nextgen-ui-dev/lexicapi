package app

import (
	"io"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var HttpLogger zerolog.Logger

func ConfigureLogger(c Config) {
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", os.ModePerm)
	}

	logFile, err := os.OpenFile("log/log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		stdlog.Fatal("Failed to open log file:", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.MessageFieldName = "msg"

	var stdout io.Writer = os.Stdout
	if c.Env == "local" {
		stdout = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	mw := zerolog.MultiLevelWriter(stdout, logFile)
	logger := zerolog.New(mw).With().Timestamp().Caller().Stack().Logger()

	log.Logger = logger

	HttpLogger = zerolog.New(mw).With().Timestamp().Logger()
}

func ReqLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		start := time.Now()
		defer func() {
			HttpLogger.Info().Fields(map[string]interface{}{
				"method":     r.Method,
				"version":    r.Proto,
				"status":     ww.Status(),
				"origin":     r.RemoteAddr,
				"host":       r.Host,
				"path":       r.URL.Path,
				"user_agent": r.Header.Get("User-Agent"),
				"latency_ms": time.Since(start).Nanoseconds() / 1000000.0,
			}).Msg("Request")
		}()

		next.ServeHTTP(ww, r)
	})
}
