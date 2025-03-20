run:
	@echo "Запуск сервиса:"
	@docker compose up -d

stop:
	@echo "Остановка сервиса:"
	@docker compose down

fmt:
	@go fmt ./...

vet:
	@go vet ./...

unit-tests: vet
	@echo "Запуск unit-тестов для основной логики сервиса:"
	@go test -v ./internal/services/...

	@echo "Запуск unit-тестов для обработчиков:"
	@go test -v ./internal/handlers/...

test-cover:
	@go test -cover ./...

clean:
	@go clean -testcache