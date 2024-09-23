# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go


# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi


# Test the application
test:
	@echo "Testing..."
	@go test ./... -v


# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@air

create-migrate:
	@echo "Create migrate seq..."
	@migrate create -seq -ext=.sql -dir=./migrations $(name)
migrate-up:
	@migrate -database postgres://tientn:Abc12345@192.168.3.8:5432/blueprint?sslmode=disable -path migrations up
migrate-down:
	@migrate -database postgres://tientn:Abc12345@192.168.3.8:5432/blueprint?sslmode=disable -path migrations down 1
migrate-force:
	@migrate -database postgres://tientn:Abc12345@192.168.3.8:5432/blueprint?sslmode=disable -path migrations force $(num)

.PHONY: all build run test clean watch create-migrate
