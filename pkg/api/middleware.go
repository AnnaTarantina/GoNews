package api

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "request_id"

// generateRequestID создает случайную строку длиной 8 символов
func generateRequestID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RequestIDMiddleware извлекает или генерирует request_id и сохраняет его в контексте
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.URL.Query().Get("request_id")
		if reqID == "" {
			reqID = generateRequestID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// responseWriter обертка для перехвата HTTP-кода ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует информацию о запросе после его выполнения
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID, _ := r.Context().Value(requestIDKey).(string)

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		log.Printf("time=%s method=%s url=%s ip=%s status=%d request_id=%s",
			start.Format(time.RFC3339),
			r.Method,
			r.URL.String(),
			r.RemoteAddr,
			rw.statusCode,
			reqID,
		)
	})
}
