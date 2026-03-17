all: test build

build:
	@echo "Building the project..."
	@mkdir -p bin
	@go build -o bin/ .

test:
	@echo "Running tests..."
	@go test -v ./...