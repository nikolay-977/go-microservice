package metrics

import (
	"net/http"
	"time"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TotalRequests - общее количество HTTP запросов
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	
	// RequestDuration - продолжительность запросов
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	
	// ActiveRequests - активные запросы
	ActiveRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of active HTTP requests",
		},
	)
	
	// ErrorRequests - ошибочные запросы
	ErrorRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_error_total",
			Help: "Total number of HTTP error requests",
		},
		[]string{"method", "endpoint", "status"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ActiveRequests)
	prometheus.MustRegister(ErrorRequests)
}

// Handler возвращает обработчик для метрик Prometheus
func Handler() http.Handler {
	return promhttp.Handler()
}

// MetricsMiddleware - middleware для сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ActiveRequests.Inc()
		
		// Перехватываем response writer для получения статуса
		rw := &responseWriter{w, http.StatusOK}
		
		next.ServeHTTP(rw, r)
		
		duration := time.Since(start).Seconds()
		ActiveRequests.Dec()
		
		// Регистрируем метрики
		TotalRequests.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Inc()
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		
		// Регистрируем ошибки (4xx и 5xx)
		if rw.status >= 400 {
			ErrorRequests.WithLabelValues(r.Method, r.URL.Path, http.StatusText(rw.status)).Inc()
		}
	})
}

// Кастомный response writer для отслеживания статуса
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
