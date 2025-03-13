# Gnose Notification Microservice

## Overview

**Gnose Notification** is a notification microservice built in **Go** that consumes messages from Redis Streams. It supports both **REST** and **WebSocket** protocols to push notifications to clients in real-time. The service enables efficient notification distribution with real-time updates for various systems using Redis as a message broker.

Ps: Im not a go developer, and this was my first experience with this language. Clearly, i used a lot of Ai agents, that let the code fully deorganized, but anyway i wanted to share my learning path.

## Architecture and Patterns

This project follows a **microservice architecture** and uses various design patterns to ensure modularity and scalability:

- **Controllers** – Handle HTTP requests and define RESTful endpoints.
- **Middleware** – Intercepts HTTP requests for logging, validation, or other cross-cutting concerns.
- **Models** – Define data structures and interfaces for communication.
- **Service Layer** – Implements the core business logic and message consumption.
- **Repository** – Manages Redis stream connections and interacts with Redis.
- **WebSocket** – Handles persistent client connections for real-time notifications.

### Design Patterns Used

- **Repository Pattern** – Encapsulates Redis stream access logic, providing a cleaner interface for service operations.
- **Service Layer Pattern** – Keeps business logic separate from the API layer for better maintainability.
- **Observer Pattern** – Used to broadcast messages to connected clients over WebSockets.
- **Singleton Pattern** – Ensures a single instance of the Redis connection and other services.

## Running the Service

To run the Gnose Notification microservice, follow these steps:

1. Clone the repository:
    ```bash
    git clone https://github.com/Davi0805/gnose-notification.git
    cd gnose-notification
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Start the service:
    ```bash
    go run main.go
    ```

- By default, the service listens on `localhost:8080`.

## Technologies and Restrictions

- **Language:** Go
- **Redis:** Used for message streaming via Redis Streams.
- **WebSocket:** Real-time communication with clients.
- **PostgreSQL:** Not explicitly required in the repo but can be used for persisting notification records.
- **Docker (optional):** Use Docker to containerize and deploy the service.

## Next Steps

- Implement user authentication and authorization.
- Improve WebSocket message handling.
- Enhance logging and error management.
- Add unit and integration tests.

## References

- [Redis Streams Documentation](https://redis.io/docs/manual/pubsub/)
- [Go Documentation](https://golang.org/doc/)
