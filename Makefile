all: test build

build:
	@echo "Building the project..."
	@go build -o bin/ .

test:
	@echo "Running tests..."
	@go test -v ./...