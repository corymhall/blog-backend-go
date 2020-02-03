package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/go-chi/chi/middleware"
)


type Middleware func(http.Handler) http.Handler


// this function creates a new logger by passing in the StructureLogger
// type to the RequestLogger function of the chi middleware package.
// the RequestLogger function expects a struct that implements the LogFormatter
// interface. In this case the StructureLogger type does implement the LogFormatter
// interface because it has the required methods, i.e. NewLogEntry. This allows us to
// use the functionality of the chi middleware package while providing our own logger.
func NewLogger(logger zerolog.Logger) Middleware {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

// StructureLogger type that implements the LogFormatter interface from the middleware package.
// It implements the interface because it has the method NewLogEntry attached.
type StructuredLogger struct {
	Logger zerolog.Logger
}

// NewLogEntry is a method of StructuredLogger. This returns a LogEntry interface.
// In this case it's the logEntry type because implements the LogEntry interface by having the
// methods Write and Panic. We do this because it allows us to provide our own logEntry type which
// contains our own custom logger, i.e. *StructureLogger rather than using the default logger provided
// by the chi middleware package. We are able to do this since the LogEntry is defined as an interface
// in the middleware package.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	
	ts := time.Now().UTC().Format(time.RFC1123)

	var req_id string
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		req_id = reqID
	}


	logger := l.Logger.With().
		Str("ts", ts).
		Str("req_id", req_id).
		Str("http_scheme", "http").
		Str("http_proto", r.Proto).
		Str("http_method", r.Method).
		Str("remote_addr", r.RemoteAddr).
		Str("user_agent", r.UserAgent()).
		Str("uri", fmt.Sprintf("%s://%s%s", "http", r.Host, r.RequestURI)).
		Logger()

	logger.Info().Msg("request started")

	entry := &StructuredLoggerEntry{logger}

	return entry
}

// type logEntry that implements the LogEntry interface of the 
// chi middleware package. It implements the interface because
// of the Write and Panic methods attacked.
type StructuredLoggerEntry struct {
	Logger zerolog.Logger
}


func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger.Info().
		Int("resp_status", status).
		Int("resp_bytes_length", bytes).
		Dur("resp_elapsed_ms", elapsed).
		Msg("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.Info().
		Str("stack", string(stack)).
		Str("panic", fmt.Sprintf("%+v", v))
}


// helper method to get the request-scoped logger
// entry.
// This allows us to do something like
// logger := Logger(r)
func Logger(r *http.Request) zerolog.Logger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}
