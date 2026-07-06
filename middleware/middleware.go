package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("HTTP запрос",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("ip", r.RemoteAddr),
			slog.String("duration", time.Since(start).String()),
		)
	})
}

// задел на будущее, если код будет более сложный и будет возможна паника
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Критический сбой обработчика запроса перехвачен (panic)",
					slog.Any("error", err),
					slog.String("path", r.URL.Path),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"success":false,"message":"Внутренняя ошибка сервера"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
