# Customer Service Application

A microservice-based application for handling customer orders and notifications.

## Setup Instructions

### 1. Configure Local Hosts

Add the following entries to your `/etc/hosts` file:

```
127.0.0.1 internal-mongodb1
127.0.0.1 internal-mongodb2
127.0.0.1 internal-mongodb3
```

### 2. Generate MongoDB Key

Generate a key file for MongoDB authentication:

```bash
openssl rand -base64 756 > db/mongo/keyfile
chmod 400 db/mongo/keyfile
```

### 3. Start Infrastructure Services

Start all required infrastructure services using Docker Compose:

```bash
# Start all services (MongoDB, Kafka, Redis, Nginx)
docker-compose up -d

# To view Kafka UI dashboard
open http://localhost:8080

# To verify Nginx is running
curl http://localhost:80
```

## Running the Application

### 1. Backend Services

```bash
# Run the customer service
cd customer
make run-with-env

# Run the order service (in a different terminal)
cd order
make run-with-env
```

### 2. Frontend Development

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start development server
npm run dev
```

Access the application via Nginx:

```
http://localhost
```

### WebSocket Connection

WebSocket server runs on port 8382:

```
ws://localhost:8382/notifications
```

### Kafka Management

Access Kafka UI dashboard:

```
http://localhost:8080
```

Monitor:

-   Topics
-   Consumer groups
-   Messages
-   Brokers status

## Architecture

The application follows a mixed architecture combining principles from both Clean Architecture and Hexagonal Architecture (Ports and Adapters).

### Core Architectural Principles

1. **Dependency Rule**: Dependencies only point inward. Inner layers don't know about outer layers.
2. **Separation of Concerns**: Business logic is separated from external concerns.
3. **Interface Segregation**: External systems are accessed through well-defined interfaces (ports).
4. **Dependency Inversion**: High-level modules don't depend on low-level modules. Both depend on abstractions.

### Layer Structure

#### 1. Domain Layer (`/domain`)

-   Core business entities and rules
-   No dependencies on other layers
-   Contains:
    -   Entities (core business objects)
    -   Value objects
    -   Domain events
    -   Domain interfaces

#### 2. Application Layer (`/application`)

-   Application-specific business rules
-   Orchestrates domain objects
-   Contains:
    -   Use cases (business operations)
    -   Ports (interfaces for external systems)
    -   DTOs and mappers
    -   Application services

#### 3. Interface Adapters

-   **API Layer** (`/api`)
    -   REST API handlers
    -   Request/Response models
    -   Route definitions
-   **Service Layer** (`/service`)
    -   Implementation of domain service interfaces
    -   Business logic implementation
    -   Orchestration of domain objects
    -   Interacts with repositories and domain objects
-   **Repository Layer** (`/repository`)
    -   Database implementations
    -   Data access patterns
    -   Query implementations
    -   Persistence layer adapters for domain entities
    -   Implements repository interfaces defined in domain layer

#### 4. Infrastructure Layer (`/infrastructure`)

-   External systems implementations
-   Framework integrations
-   Contains:
    -   Kafka implementation
    -   Database connections
    -   External service clients

#### 5. Supporting Components

-   **Component Layer** (`/component`)
    -   Application context
    -   Dependency injection
    -   Shared components
-   **Config Layer** (`/config`)
    -   Configuration management
    -   Environment settings
-   **Middleware** (`/middleware`)
    -   HTTP middleware
    -   Request processing
-   **WebSocket** (`/websocket`)
    -   Real-time communication
    -   Connection management

### Ports and Adapters Pattern

#### Ports (Interfaces)

-   Located in application layer under `application/port` package
-   Define how the application interacts with external systems through clean interfaces
-   Key message ports:
    -   `MessageProducer`: Interface for sending messages to message brokers
        -   `Publish(message domain.Message) error`: Publishes domain messages
        -   `Close() error`: Cleanly shuts down the producer
    -   `MessageConsumer`: Interface for receiving messages from message brokers
        -   `RegisterHandler(topic string, handler func(domain.Message) error) error`: Registers message handlers
        -   `Start(ctx context.Context) error`: Starts consuming messages
        -   `Close() error`: Cleanly shuts down the consumer
    -   `DatabasePort`: Interface for database operations
        -   `Query(ctx, collection, filter, result)`: Type-safe query for multiple documents
        -   `QueryOne(ctx, collection, filter, result)`: Type-safe query for single document
        -   `Insert`: Inserts documents
        -   `Update`: Updates documents
        -   `Delete`: Deletes documents
        -   `Upsert`: Updates or inserts documents
        -   `Incr`: Atomically increments numeric fields
        -   `FindOneAndUpdate`: Atomic update with return value
        -   `Close`: Auto-closes with parent context

#### Adapters (Implementations)

-   Located in infrastructure layer under specific technology packages:
    -   `infrastructure/kafka`: Kafka message broker implementations
    -   `infrastructure/mongodb`: MongoDB database implementations
-   Key adapters:
    -   Kafka adapters:
        -   `Producer`: Implements `MessageProducer` interface for Kafka
        -   `Consumer` and `RetryableConsumer`: Implement `MessageConsumer` interface for Kafka
        -   Provides additional Kafka-specific functionality like:
            -   Admin operations (topic management)
            -   Configuration management
            -   Retry mechanisms for reliable message processing
    -   MongoDB adapters:
        -   `MongoAdapter`: Implements `DatabasePort` interface
        -   Features:
            -   Type-safe query operations with generics
            -   Automatic connection management with context
            -   Built-in error wrapping with context
            -   Support for read/write separation
            -   Atomic operations for ID generation

The Ports and Adapters pattern in this project ensures:

-   Business logic remains independent of messaging implementation details
-   Easy switching between different message broker implementations if needed
-   Clear separation between domain logic and infrastructure concerns
-   Testability through interface mocking
-   Automatic resource cleanup through context management

### Repository Implementation

#### MongoDB Repository Pattern

-   Repositories use the MongoDB adapter through the DatabasePort interface
-   Support for read/write separation:

    ```go
    type MongoRepository struct {
        writeDB *mongodb.MongoAdapter  // For write operations
        readDB  *mongodb.MongoAdapter  // For read operations (optional)
    }
    ```

-   Type-safe query operations:

    ```go
    // Example: Query multiple documents
    var orders []Order
    err := repo.readDB.Query(ctx, collection, filter, &orders)

    // Example: Query single document
    var order Order
    err := repo.readDB.QueryOne(ctx, collection, filter, &order)
    ```

-   Atomic operations (e.g., ID generation):

    ```go
    var idGen IDGen
    err := repo.writeDB.FindOneAndUpdate(
        ctx, collection,
        bson.M{"key": key},
        bson.M{"$inc": bson.M{"value": 1}},
        &idGen,
        options.FindOneAndUpdate().SetUpsert(true),
    )
    ```

-   Automatic connection management:
    ```go
    adapter, err := mongodb.NewMongoAdapter(ctx, uri, dbName)
    // Connection will be automatically closed when ctx is done
    ```

### Technology Stack

#### Backend

-   Go (Golang)
-   MongoDB (Replica Set)
-   Kafka for message streaming
-   Redis for caching
-   WebSockets for real-time communication

#### Frontend

-   Next.js
-   TypeScript
-   Tailwind CSS
-   React Hot Toast

### Service vs. Use Case Layer

#### Service Layer Rules

-   **Domain-Specific**: Each service is dedicated to a single domain entity (e.g., `OrderService` for the Order domain)
-   **Optional Layer**: May be omitted in simpler applications where domain logic is minimal
-   **Implementation Rules**:
    -   Implements domain service interfaces defined in the domain layer
    -   Contains domain-specific business logic and validation
    -   Typically works with a single repository type
    -   Should not orchestrate multiple domain services

#### Use Case Layer Rules

-   **Cross-Domain**: Orchestrates operations across multiple domains (e.g., `CreateOrderUseCase` might involve Order, Customer, and Payment domains)
-   **Higher Abstraction**: Implements application-specific business flows
-   **Implementation Rules**:
    -   Located in the application layer
    -   May coordinate multiple domain services
    -   Handles transaction boundaries
    -   Implements application-level validation and business rules
    -   Performs data transformations between the domain and external layers

When to use which:

-   Use **Service Layer** when the logic is specific to a single domain entity
-   Use **Use Case Layer** when orchestrating operations across multiple domains
-   Complex applications typically have both layers
-   Simpler applications might only need one of these layers
