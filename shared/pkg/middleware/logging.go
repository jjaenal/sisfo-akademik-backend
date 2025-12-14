package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &respWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		dur := time.Since(start)
		l.Info("request", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Int("status", ww.status), zap.Duration("latency", dur))
	})
}

type respWriter struct {
	http.ResponseWriter
	status int
}

func (w *respWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
