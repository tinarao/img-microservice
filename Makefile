# Переменные
BINARY_NAME=api
MAIN_FILE=cmd/main.go

run: build
	go run $(MAIN_FILE)

build: deps
	go build -o $(BINARY_NAME) $(MAIN_FILE)

deps:
	go mod tidy
	
test:
	go test ./... -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
