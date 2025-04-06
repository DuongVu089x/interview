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

## Architecture

### Backend

The backend is organized using a clean architecture approach:

-   **Application Layer**

    -   Contains ports (interfaces) for communicating with infrastructure services like Kafka and message brokers. Includes interfaces for MessageProducer and MessageConsumer
    -   Containt domain layer for business logic
        -   dto that defines the data transfer object
        -   mapper that defines the mapper for the data transfer object
        -   usercase that defines the business logic

-   **Component Layer**

    -   Includes `appctx` for application context management

-   **Domain Layer**

    -   Contains core business entities and logic for domain including
        -   entity that defines the structure of the data
        -   repository that defines the operations on the data
        -   service that defines the business logic

-   **Infrastructure Layer**

    -   Implements interfaces defined in the application layer
    -   Connects to external services (Kafka, Redis, MongoDB)

-   **Repository Layer**

    -   Implements interfaces defined in the repository layer, with the implementation SQL queries database
    -   Connects to external services (MongoDB, PostgreSQL, etc)

-   **Service Layer**

    -   Implements interfaces defined in the service layer, with the implementation of the sample business logic

-   **Main Application**

    -   Entry point that wires everything together

-   **WebSocket**
    -   Handles real-time notifications
    -   Dispatch notification to the client, contains the logic business for dispatching notification to the client
    -   Handler that is custom for application to adapt the business logic to the websocket
    -   Common file for websocket server, router, connection. It should be move to sdk

### Frontend (Web Client)

A Next.js application for real-time notifications:

-   Connects to WebSocket server for real-time updates
-   Displays notifications with toast alerts
-   Provides a dashboard interface

## Technology Stack

-   **Backend**

    -   Go (Golang)
    -   MongoDB (Replica Set)
    -   Kafka for message streaming
    -   Redis for caching
    -   WebSockets for real-time communication

-   **Frontend**
    -   Next.js
    -   TypeScript
    -   Tailwind CSS
    -   React Hot Toast

## Running the Application

### Backend

```bash
# Run the customer service
cd customer
make run-with-env
```

```bash
# Run the order service
cd order
make run-with-env
```

### Frontend

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start development server
npm run dev
```

### WebSocket Connection

The WebSocket server runs on port 8382. The frontend automatically connects to:

```
ws://localhost:8382/notifications
```

## Testing

A test HTML page is provided at `/web/public/test.html` for testing the WebSocket connection manually.
