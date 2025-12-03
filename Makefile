.PHONY: build run test clean docker-build docker-run load-test

build:
	go mod tidy
	go build -o bin/main main.go

run: build
	./bin/main

clean:
	rm -rf bin/
	rm -f audit.log

docker-build:
	docker build -t go-microservice:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

load-test:
	@echo "Запуск нагрузочного тестирования..."
	wrk -t12 -c500 -d60s http://localhost:8080/api/users

health-check:
	curl http://localhost:8080/health

metrics:
	curl http://localhost:8080/metrics | head -50
