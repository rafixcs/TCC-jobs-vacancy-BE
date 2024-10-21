package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)

		lrw := &LoggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)

		log.Printf("Completed %s %s in %v with status %d", r.Method, r.RequestURI, time.Since(start), lrw.statusCode)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
