package utils

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/time/rate"
)

// Default значения на случай отсутствия переменных окружения
const (
	DefaultRateLimit  = 1000
	DefaultBurstLimit = 5000
)

// Создаем глобальный rate limiter
var limiter *rate.Limiter

// Инициализация rate limiter при импорте пакета
func init() {
	// Получаем значения из переменных окружения или используем значения по умолчанию
	rateLimit := getEnvAsInt("RATE_LIMIT_PER_SECOND", DefaultRateLimit)
	burstLimit := getEnvAsInt("RATE_LIMIT_BURST", DefaultBurstLimit)

	fmt.Printf("Инициализация знчений: rate limiter = %d запросов в секунду, burst = %d\n", rateLimit, burstLimit)

	limiter = rate.NewLimiter(rate.Limit(rateLimit), burstLimit)
}

// getEnvAsInt получает значение переменной окружения как целое число
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		// В случае ошибки парсинга возвращаем значение по умолчанию
		return defaultValue
	}

	return value
}

// RateLimitMiddleware - middleware для ограничения скорости запросов
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
