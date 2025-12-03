package services

import (
	"fmt"
	"sync"
)

type IntegrationService struct {
	notifications chan string
	wg            sync.WaitGroup
	mu            sync.RWMutex
}

var (
	integrationInstance *IntegrationService
	integrationOnce     sync.Once
)

// GetIntegrationService возвращает синглтон сервиса интеграций
func GetIntegrationService() *IntegrationService {
	integrationOnce.Do(func() {
		integrationInstance = &IntegrationService{
			notifications: make(chan string, 1000),
		}
		integrationInstance.startWorkerPool(5) // 5 воркеров для обработки уведомлений
	})
	return integrationInstance
}

// startWorkerPool запускает пул воркеров для обработки уведомлений
func (s *IntegrationService) startWorkerPool(workerCount int) {
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.notificationWorker(i + 1)
	}
}

// notificationWorker обрабатывает уведомления
func (s *IntegrationService) notificationWorker(id int) {
	defer s.wg.Done()

	for notification := range s.notifications {
		// Имитация отправки уведомления
		fmt.Printf("Worker %d: Отправка уведомления: %s\n", id, notification)
	}
}

// SendNotification отправляет асинхронное уведомление
func (s *IntegrationService) SendNotification(message string) {
	select {
	case s.notifications <- message:
		// Уведомление отправлено в канал
	default:
		fmt.Println("Очередь уведомлений переполнена")
	}
}

// Shutdown корректно останавливает сервис
func (s *IntegrationService) Shutdown() {
	close(s.notifications)
	s.wg.Wait()
	fmt.Println("Сервис интеграций остановлен")
}
