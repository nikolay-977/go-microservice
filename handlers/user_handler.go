package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
	
	"go-microservice/models"
	"go-microservice/services"
	"go-microservice/utils"
)

// GetAllUsersHandler возвращает всех пользователей
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	userService := services.GetUserService()
	users := userService.GetAllUsers()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserHandler возвращает пользователя по ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}
	
	userService := services.GetUserService()
	user, err := userService.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUserHandler создает нового пользователя
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Имя и email обязательны", http.StatusBadRequest)
		return
	}
	
	userService := services.GetUserService()
	createdUser := userService.CreateUser(user)
	
	// Асинхронное логирование
	go utils.LogUserAction("CREATE", createdUser.ID)
	
	// Асинхронная отправка уведомления
	integrationService := services.GetIntegrationService()
	go integrationService.SendNotification("Создан новый пользователь: " + createdUser.Name)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// UpdateUserHandler обновляет пользователя
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}
	
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	if user.Name == "" || user.Email == "" {
		http.Error(w, "Имя и email обязательны", http.StatusBadRequest)
		return
	}
	
	userService := services.GetUserService()
	updatedUser, err := userService.UpdateUser(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Асинхронное логирование
	go utils.LogUserAction("UPDATE", updatedUser.ID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUserHandler удаляет пользователя
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}
	
	userService := services.GetUserService()
	err = userService.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Асинхронное логирование
	go utils.LogUserAction("DELETE", id)
	
	// Асинхронная отправка уведомления
	integrationService := services.GetIntegrationService()
	go integrationService.SendNotification("Удален пользователь с ID: " + strconv.Itoa(id))
	
	w.WriteHeader(http.StatusNoContent)
}
