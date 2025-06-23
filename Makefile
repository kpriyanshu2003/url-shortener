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