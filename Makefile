generate-swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/main.go --output ./docs