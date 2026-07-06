package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fit24/handler"
	"fit24/middleware"
	"fit24/repository"
	"fit24/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	logLevelStr := getEnv("LOG_LEVEL", "info")
	var level slog.Level
	switch strings.ToLower(logLevelStr) {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})))

	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "fit24_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	slog.Info("успешное подключение к БД")

	dbTimeout := getDurationEnv("DB_TIMEOUT", "3s")
	readTimeout := getDurationEnv("SERVER_READ_TIMEOUT", "10s")
	writeTimeout := getDurationEnv("SERVER_WRITE_TIMEOUT", "10s")
	shutdownTimeout := getDurationEnv("SHUTDOWN_TIMEOUT", "5s")

	repo := repository.NewPostgresRepository(pool)
	leadService := service.NewLeadService(repo)
	httpHandler := handler.NewHTTPHandler(leadService, dbTimeout)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/order", httpHandler.HandleOrder)
	mux.HandleFunc("/api/contact", httpHandler.HandleContact)
	mux.HandleFunc("/api/health", httpHandler.HandleHealth)

	finalHandler := middleware.Recovery(middleware.Logging(mux))

	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      finalHandler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	go func() {
		slog.Info("сервер запущен", slog.String("url", "http://localhost:"+port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	slog.Info("получен сигнал плавной остановки сервера")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("ошибка при принудительной остановке сервера", slog.String("error", err.Error()))
	} else {
		slog.Info("все HTTP-соединения успешно завершены")
	}

	slog.Info("сервер полностью остановлен")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDurationEnv(key, fallback string) time.Duration {
	val := getEnv(key, fallback)
	d, err := time.ParseDuration(val)
	if err != nil {
		return time.Duration(0)
	}
	return d
}
