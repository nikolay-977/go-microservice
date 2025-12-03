package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// Создаем rate limiter: 80000 запросов в секунду с burst 400000
var limiter = rate.NewLimiter(rate.Limit(80000), 400000)

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
