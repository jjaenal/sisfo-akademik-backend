package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
			w.Header().Set("X-Request-ID", id)
			r.Header.Set("X-Request-ID", id)
		} else {
			w.Header().Set("X-Request-ID", id)
			r.Header.Set("X-Request-ID", id)
		}
		next.ServeHTTP(w, r)
	})
}

func Logging(l *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &respWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(ww, r)
		dur := time.Since(start)

		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("request_id", w.Header().Get("X-Request-ID")),
			zap.Int("status", ww.status),
			zap.Duration("latency", dur),
		}

		// Extract TraceID and SpanID from context
		spanCtx := trace.SpanContextFromContext(r.Context())
		if spanCtx.HasTraceID() {
			fields = append(fields, zap.String("trace_id", spanCtx.TraceID().String()))
		}
		if spanCtx.HasSpanID() {
			fields = append(fields, zap.String("span_id", spanCtx.SpanID().String()))
		}

		l.Info("request", fields...)
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
