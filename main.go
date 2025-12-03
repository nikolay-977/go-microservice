package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/gorilla/mux"
	"go-microservice/handlers"
	"go-microservice/metrics"
	"go-microservice/utils"
)

func main() {
	// Инициализируем роутер
	r := mux.NewRouter()
	
	// Настройка middleware
	r.Use(metrics.MetricsMiddleware)
	r.Use(utils.RateLimitMiddleware)
	
	// Регистрируем маршруты
	handlers.RegisterRoutes(r)
	
	// Маршрут для метрик Prometheus
	r.Handle("/metrics", metrics.Handler())
	
	// Настройка HTTP сервера
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	fmt.Println("Сервер запущен на порту :8080")
	log.Fatal(srv.ListenAndServe())
}
