generate-swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/main.go --output ./docs

run-dev:
	@echo "Running application in development mode..."
	air .

build: 
	@echo "Building the application..."
	go build -o bin/app cmd/main.go

run: generate-swagger build
	@echo "Running the application..."
	./bin/app

test-all:
	@echo "Running all tests..."
	go test -v ./...

test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-api:
	@echo "Running API tests..."
	go test -v ./tests/api/...


coverage:
	@echo "Generating test coverage report..."
	go test -v -coverpkg=./internal/database,./internal/controller -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out
	@echo "Coverage report generated: coverage.html"

coverage-view:
	@echo "Opening coverage report..."
	go tool cover -func=coverage.out

build-docker:
	@echo "Building Docker image..."
	COMPOSE_BAKE=true docker build -t url-shortener:latest .