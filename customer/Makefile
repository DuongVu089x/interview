# Makefile for Order Service

# Environment Variables (can be overridden via environment)
MONGODB_URI ?= mongodb://root:example@localhost:27020,localhost:27018,localhost:27019/?replicaSet=rs0&authSource=admin
MONGODB_READ_URI ?= mongodb://root:example@localhost:27020,localhost:27018,localhost:27019/?replicaSet=rs0&authSource=admin&readPreference=secondary
KAFKA_BOOTSTRAP_SERVERS ?= localhost:9092
KAFKA_SECURITY_PROTOCOL ?= plaintext
KAFKA_DEFAULT_TOPIC ?= orders-topic
REDIS_ADDR ?= localhost:6379
REDIS_PASSWORD ?=
REDIS_DB ?= 0
SERVER_PORT ?= 8082
GRPC_PORT ?= 50051

# Jaeger Configuration
JAEGER_ENABLED ?= true
JAEGER_AGENT_HOST ?= 0.0.0.0
JAEGER_AGENT_PORT ?= 14268
JAEGER_SERVICE_NAME ?= customer-service

# Export all variables for child processes
export MONGODB_URI
export MONGODB_READ_URI
export KAFKA_BOOTSTRAP_SERVERS
export KAFKA_SECURITY_PROTOCOL
export KAFKA_DEFAULT_TOPIC
export REDIS_ADDR
export REDIS_PASSWORD
export REDIS_DB
export SERVER_PORT
export GRPC_PORT
export JAEGER_ENABLED
export JAEGER_AGENT_HOST
export JAEGER_AGENT_PORT
export JAEGER_SERVICE_NAME

# Run the application
run:
	@echo "Running order service..."
	@go run main.go

# Run with environment variables from Makefile
run-with-env:
	@echo "Running order service with environment variables..."
	@MONGODB_URI=$(MONGODB_URI) \
	MONGODB_READ_URI=$(MONGODB_READ_URI) \
	KAFKA_BOOTSTRAP_SERVERS=$(KAFKA_BOOTSTRAP_SERVERS) \
	KAFKA_SECURITY_PROTOCOL=$(KAFKA_SECURITY_PROTOCOL) \
	KAFKA_DEFAULT_TOPIC=$(KAFKA_DEFAULT_TOPIC) \
	REDIS_ADDR=$(REDIS_ADDR) \
	REDIS_PASSWORD=$(REDIS_PASSWORD) \
	REDIS_DB=$(REDIS_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	GRPC_PORT=$(GRPC_PORT) \
	JAEGER_ENABLED=$(JAEGER_ENABLED) \
	JAEGER_AGENT_HOST=$(JAEGER_AGENT_HOST) \
	JAEGER_AGENT_PORT=$(JAEGER_AGENT_PORT) \
	JAEGER_SERVICE_NAME=$(JAEGER_SERVICE_NAME) \
	go run main.go

# Build the application
build:
	@echo "Building order service..."
	@go build -o order-service ./

# Build gRPC test client
build-grpc-client:
	go build -o bin/grpc-client cmd/grpc-client/main.go

# Run gRPC test client with a test user ID
test-grpc:
	@if [ -z "$(user)" ]; then \
		echo "Usage: make test-grpc user=<user_id>"; \
		exit 1; \
	fi
	go run cmd/grpc-client/main.go -user=$(user)

.PHONY: run run-with-env build build-grpc-client test-grpc
