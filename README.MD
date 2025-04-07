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
-   **Repository Layer** (`/repository`)
    -   Database implementations
    -   Data access patterns
    -   Query implementations

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

-   Located in application layer with `_port` suffix
-   Define how the application interacts with external systems
-   Example: `messaging_port` for Kafka interaction

#### Adapters (Implementations)

-   Located in infrastructure and repository layers
-   Implement the port interfaces
-   Connect to external systems
-   Example: Kafka consumer/producer implementations

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

Access options:

1. Direct: `http://localhost:3000`
2. Via Nginx: `http://localhost`

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

### Testing

A test HTML page is available at `/web/public/test.html` for testing WebSocket connections manually.
