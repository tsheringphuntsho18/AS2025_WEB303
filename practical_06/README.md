# Practical_06 Report: Comprehensive Testing for Microservices

This project implements a Student Cafe management system using a microservices architecture. The services are written in Go and communicate with each other using gRPC. An API Gateway is provided to expose a RESTful HTTP interface to external clients. The system includes a comprehensive testing suite covering unit, integration, and end-to-end (E2E) tests.

## Architecture

The system is composed of the following services, each running in its own Docker container:

-   **API Gateway**: The single entry point for all external clients. It receives HTTP requests and translates them into gRPC calls to the appropriate backend microservice.
-   **User Service**: Manages user accounts (students, cafe owners). It exposes a gRPC interface for CRUD operations.
-   **Menu Service**: Manages cafe menu items and their prices. It exposes a gRPC interface.
-   **Order Service**: Manages customer orders. This service communicates with the User and Menu services via gRPC to validate user existence and to snapshot menu item prices at the time of order.
-   **Databases**: Each microservice has its own dedicated PostgreSQL database (`user-db`, `menu-db`, `order-db`) to ensure loose coupling and data isolation.
-   **Protobufs (`student-cafe-protos`)**: A centralized repository containing the Protocol Buffer definitions (`.proto`) for all gRPC services, ensuring a single source of truth for the service contracts.

### Service Communication Flow
- **Client → API Gateway (HTTP REST)**: A client (e.g., a web browser or `curl`) sends an HTTP request to the API Gateway.
- **API Gateway → Microservices (gRPC)**: The gateway translates the HTTP request into a gRPC call and forwards it to the corresponding internal service (User, Menu, or Order).
- **Order Service → User/Menu Services (gRPC)**: When creating an order, the Order Service makes gRPC calls to the User Service to verify the user exists and to the Menu Service to get the current price of items.

## Features

-   **Microservice Architecture**: Decoupled services for improved scalability and maintainability.
-   **gRPC Communication**: High-performance, strongly-typed RPC framework for inter-service communication.
-   **Centralized Protobufs**: A dedicated module (`student-cafe-protos`) for managing gRPC service definitions.
-   **API Gateway Pattern**: A single, unified REST API exposed to the outside world.
-   **Containerized**: Fully containerized with Docker and orchestrated using Docker Compose.
-   **Comprehensive Testing**:
    -   **Unit Tests**: Test individual gRPC server methods in isolation using an in-memory SQLite database.
    -   **Integration Tests**: Test the gRPC communication and logic between services using in-memory connections (`bufconn`).
    -   **End-to-End (E2E) Tests**: Test the full application flow by making HTTP requests to the running services via the API Gateway.
-   **Makefile Automation**: Streamlined commands for building, testing, running, and cleaning the environment.

## Getting Started

### Prerequisites

-   [Go](https://go.dev/doc/install) (version 1.22 or later)
-   [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
-   [Protocol Buffers Compiler (`protoc`)](https://grpc.io/docs/protoc-installation/)
-   Go gRPC plugins (`make install-tools` will install these for you)

### Installation and Running

You can set up and run the entire application using either the `Makefile` or the deployment script.

**1. Using the Makefile (Recommended)**

This two-step process installs dependencies, generates Protobuf code, builds Docker images, and starts all services.

```bash
# First, run the one-time development setup
make dev-setup

# Then, start all services
make docker-up
```

**2. Using the Deployment Script**

This single script handles building, deploying, and starting all services.

```bash
./deploy.sh
```

After starting the services, the API Gateway will be accessible at `http://localhost:8080`.

## Usage

All API endpoints are accessed through the API Gateway at `http://localhost:8080`.

### API Endpoints

-   **User Service**
    -   `POST /api/users`: Create a new user.
    -   `GET /api/users`: Get a list of all users.
    -   `GET /api/users/{id}`: Get a specific user by their ID.
-   **Menu Service**
    -   `POST /api/menu`: Create a new menu item.
    -   `GET /api/menu`: Get a list of all menu items.
    -   `GET /api/menu/{id}`: Get a specific menu item by its ID.
-   **Order Service**
    -   `POST /api/orders`: Create a new order.
    -   `GET /api/orders`: Get a list of all orders.
    -   `GET /api/orders/{id}`: Get a specific order by its ID.

### Example `curl` Commands

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name": "John Doe", "email": "john.doe@example.com", "is_cafe_owner": false}'

# Create a menu item
curl -X POST http://localhost:8080/api/menu \
  -H 'Content-Type: application/json' \
  -d '{"name": "Espresso", "description": "Strong black coffee", "price": 3.00}'

# Create an order (uses user_id=1 and menu_item_id=1)
curl -X POST http://localhost:8080/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"user_id": 1, "items": [{"menu_item_id": 1, "quantity": 2}]}'

# Get all orders
curl http://localhost:8080/api/orders
```

## Testing

The project includes a multi-layered testing strategy. You can run tests using the `Makefile`.

-   **Run Unit Tests:**
    These tests check individual gRPC server functions using an in-memory SQLite database.
    ```bash
    make test-unit
    ```

-   **Run Integration Tests:**
    These tests verify the gRPC interactions between services using in-memory connections, without needing a full Docker environment.
    ```bash
    make test-integration
    ```

-   **Run End-to-End (E2E) Tests:**
    These tests spin up the entire Docker environment and validate the complete application flow through HTTP requests to the API Gateway.
    ```bash
    # This command starts the Docker containers, runs tests, and then stops them.
    make test-e2e-docker
    ```

-   **Run All Tests:**
    To run all unit, integration, and E2E tests in sequence:
    ```bash
    make test-all
    ```

-   **Generate Test Coverage:**
    This command runs unit tests and generates an HTML coverage report for each service.
    ```bash
    make test-coverage
    # Reports are generated at:
    # - user-service/coverage.html
    # - menu-service/coverage.html
    # - order-service/coverage.html
    ```

## Makefile Commands

The `Makefile` provides several useful commands for development and management.

| Command                  | Description                                                                  |
| ------------------------ | ---------------------------------------------------------------------------- |
| `make help`              | Shows a list of all available commands.                                      |
| `make dev-setup`         | Installs dependencies, generates Protobuf code, and builds Docker images.    |
| `make proto-generate`    | Generates gRPC and Protobuf Go code from `.proto` files.                     |
| `make docker-build`      | Builds or rebuilds the Docker images for all services.                       |
| `make docker-up`         | Starts all services using Docker Compose.                                    |
| `make docker-down`       | Stops and removes all running service containers.                            |
| `make docker-logs`       | Tails the logs from all running services.                                    |
| `make test-unit`         | Runs unit tests for all services.                                            |
| `make test-integration`  | Runs integration tests.                                                      |
| `make test-e2e-docker`   | Runs E2E tests, managing the Docker environment automatically.               |
| `make test-all`          | Runs all unit, integration, and E2E tests.                                   |
| `make test-coverage`     | Runs unit tests and generates HTML coverage reports.                         |
| `make clean`             | Stops containers, cleans up test artifacts, and removes Docker volumes.      |