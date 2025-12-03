package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	loggerInstance *Logger
	once           sync.Once
)

type Logger struct {
	file *os.File
	mu   sync.Mutex
}

func GetLogger() *Logger {
	once.Do(func() {
		// Создаем или открываем файл лога
		file, err := os.OpenFile("audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		
		loggerInstance = &Logger{file: file}
	})
	return loggerInstance
}

func (l *Logger) Log(action string, userID int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("[%s] Action: %s, UserID: %d\n", timestamp, action, userID)
	
	if _, err := l.file.WriteString(message); err != nil {
		log.Printf("Ошибка записи в лог: %v", err)
	}
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// Вспомогательная функция для логирования действий пользователя
func LogUserAction(action string, userID int) {
	go func() {
		logger := GetLogger()
		logger.Log(action, userID)
		fmt.Printf("Логирование: %s для пользователя %d\n", action, userID)
	}()
}
