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
SERVER_PORT ?= 8081

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
	go run main.go

# Build the application
build:
	@echo "Building order service..."
	@go build -o order-service ./

.PHONY: run run-with-env build
