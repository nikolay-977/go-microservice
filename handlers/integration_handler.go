package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// HealthCheckHandler возвращает статус health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"service": "User Management Service",
	})
}

// RegisterRoutes регистрирует все маршруты
func RegisterRoutes(r *mux.Router) {
	// API маршруты
	api := r.PathPrefix("/api").Subrouter()

	// Пользователи
	api.HandleFunc("/users", GetAllUsersHandler).Methods("GET")
	api.HandleFunc("/users/{id}", GetUserHandler).Methods("GET")
	api.HandleFunc("/users", CreateUserHandler).Methods("POST")
	api.HandleFunc("/users/{id}", UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/users/{id}", DeleteUserHandler).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
}
