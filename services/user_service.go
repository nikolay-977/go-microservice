package services

import (
	"errors"
	"sync"
	
	"go-microservice/models"
)

type UserService struct {
	users     map[int]models.User
	mu        sync.RWMutex
	currentID int
}

var (
	instance *UserService
	once     sync.Once
)

// GetUserService возвращает синглтон сервиса пользователей
func GetUserService() *UserService {
	once.Do(func() {
		instance = &UserService{
			users:     make(map[int]models.User),
			currentID: 1,
		}
		// Добавляем тестового пользователя
		instance.users[1] = models.User{ID: 1, Name: "Тестовый пользователь", Email: "test@example.com"}
	})
	return instance
}

// GetAllUsers возвращает всех пользователей
func (s *UserService) GetAllUsers() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// GetUserByID возвращает пользователя по ID
func (s *UserService) GetUserByID(id int) (models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	user, exists := s.users[id]
	if !exists {
		return models.User{}, errors.New("пользователь не найден")
	}
	return user, nil
}

// CreateUser создает нового пользователя
func (s *UserService) CreateUser(user models.User) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	user.ID = s.currentID
	s.users[s.currentID] = user
	s.currentID++
	
	return user
}

// UpdateUser обновляет существующего пользователя
func (s *UserService) UpdateUser(id int, updatedUser models.User) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.users[id]; !exists {
		return models.User{}, errors.New("пользователь не найден")
	}
	
	updatedUser.ID = id
	s.users[id] = updatedUser
	return updatedUser, nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.users[id]; !exists {
		return errors.New("пользователь не найден")
	}
	
	delete(s.users, id)
	return nil
}
