package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"fit24/handler"
	"fit24/repository"
	"fit24/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

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
		log.Fatal("Ошибка конфигурации пула pgx: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Не удалось подключиться к базе данных: ", err)
	}
	fmt.Println("Успешное подключение к БД")

	repo := repository.NewPostgresRepository(pool)
	leadService := service.NewLeadService(repo)
	httpHandler := handler.NewHTTPHandler(leadService)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/order", httpHandler.HandleOrder)
	http.HandleFunc("/api/contact", httpHandler.HandleContact)

	port := getEnv("PORT", "8080")
	fmt.Printf("Сервер успешно запущен по адресу http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
