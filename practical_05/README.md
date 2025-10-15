# Practical_05 Report: Microservices Architecture - Student Cafe System

## Overview

This practical demonstrates the transformation of a monolithic Student Cafe application into a distributed microservices architecture using Go, Docker, PostgreSQL, and Consul for service discovery.

## Architecture Diagram

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    │    (Port 8080)  │
                    └─────────┬───────┘
                              │
                    ┌─────────┴───────┐
                    │     Consul      │
                    │  Service Disc.  │
                    │   (Port 8500)   │
                    └─────────┬───────┘
                              │
            ┌─────────────────┼─────────────────┐
            │                 │                 │
    ┌───────▼──────┐ ┌───────▼──────┐ ┌───────▼──────┐
    │User Service  │ │Menu Service  │ │Order Service │
    │ (Port 8081)  │ │ (Port 8082)  │ │ (Port 8083)  │
    └───────┬──────┘ └───────┬──────┘ └───────┬──────┘
            │                │                │
    ┌───────▼──────┐ ┌───────▼──────┐ ┌───────▼──────┘
    │   User DB    │ │   Menu DB    │ │   Order DB   │
    │ (Port 5434)  │ │ (Port 5433)  │ │ (Port 5435)  │
    └──────────────┘ └──────────────┘ └──────────────┘
```

## Service Boundaries Justification

### 1. User Service (Port 8081)

**Responsibility:** User management and authentication
**Database:** `user_db` on port 5434
**Endpoints:**

- `POST /users` - Create new user
- `GET /users/{id}` - Get user by ID
- `GET /health` - Health check

**Justification:**

- Single responsibility: Manages only user-related operations
- Independent scalability: Can scale based on user registration patterns
- Security isolation: User data is isolated in its own database
- Reusability: Can be used by multiple other services

### 2. Menu Service (Port 8082)

**Responsibility:** Menu item management
**Database:** `menu_db` on port 5433
**Endpoints:**

- `POST /menu` - Create menu item
- `GET /menu/{id}` - Get menu item by ID
- `GET /health` - Health check

**Justification:**

- Business domain separation: Menu management is distinct from orders/users
- Independent updates: Menu changes don't affect user or order services
- Performance optimization: Can be cached independently
- Different scaling needs: Menu reads vs order writes have different patterns

### 3. Order Service (Port 8083)

**Responsibility:** Order processing and management
**Database:** `order_db` on port 5435
**Endpoints:**

- `POST /orders` - Create new order
- `GET /orders` - Get all orders
- `GET /health` - Health check

**Justification:**

- Complex business logic: Orders involve validation across multiple services
- Transaction management: Order creation requires coordination
- Audit trail: Orders need complete history tracking
- Integration point: Orchestrates user and menu services

### 4. API Gateway (Port 8080)

**Responsibility:** Request routing and service discovery
**Routes:**

- `/api/users/*` → User Service
- `/api/menu/*` → Menu Service
- `/api/orders/*` → Order Service

**Justification:**

- Single entry point: Simplifies client integration
- Service abstraction: Clients don't need to know individual service locations
- Load balancing: Can distribute requests across service instances
- Cross-cutting concerns: Authentication, logging, rate limiting

## Inter-Service Communication

### Order Creation Flow

1. **API Gateway** receives order request at `/api/orders`
2. **Order Service** validates request structure
3. **Order Service** → **User Service**: `GET /users/{id}` to validate user exists
4. **Order Service** → **Menu Service**: `GET /menu/{id}` for each item to validate and get current price
5. **Order Service** creates order with validated data
6. Response sent back through API Gateway

### Service Discovery

- All services register with Consul on startup
- Health checks every 10 seconds via `/health` endpoint
- API Gateway dynamically discovers service locations
- Automatic failover if service instances become unhealthy

## Challenges Encountered and Solutions

### 1. Service Discovery and Communication

**Challenge:** Services needed to communicate without hardcoded URLs
**Solution:**

- Implemented Consul for service registration and discovery
- API Gateway dynamically resolves service locations
- Health checks ensure only healthy services receive traffic

### 2. Database Independence

**Challenge:** Ensuring complete data isolation between services
**Solution:**

- Separate PostgreSQL containers for each service
- Different ports (5433, 5434, 5435) for each database
- Independent schemas and migrations per service

### 3. Inter-Service Data Validation

**Challenge:** Order service needs to validate users and menu items from other services
**Solution:**

- HTTP calls to validate data existence
- Snapshot pricing at order creation time
- Graceful error handling for service unavailability

### 4. Container Networking

**Challenge:** Services communicating within Docker network
**Solution:**

- Used service names as hostnames in Docker Compose
- Proper dependency ordering with `depends_on`
- Internal network communication on service ports

### 5. Development vs Production Configuration

**Challenge:** Different database URLs for local vs containerized development
**Solution:**

- Environment variables for database configuration
- Fallback defaults for local development
- Container-specific environment variables in docker-compose.yml

## Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Go 1.21+ (for local development)
- PostgreSQL client (optional, for database inspection)

## API Endpoints

### User Service (via API Gateway: `/api/users`)

- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user

### Menu Service (via API Gateway: `/api/menu`)

- `POST /api/menu` - Create menu item
- `GET /api/menu/{id}` - Get menu item

### Order Service (via API Gateway: `/api/orders`)

- `POST /api/orders` - Create order
- `GET /api/orders` - List orders

## Monitoring and Health Checks

### Consul Health Monitoring

- Access Consul UI at http://localhost:8500
- All services register with health checks
- 10-second intervals with 3-second timeouts
- HTTP health checks on `/health` endpoints

## Reflection Essay

### Monolith vs Microservices:
The choice between monolithic and microservices architectures for the Student Cafe system depends on scale and complexity. A monolithic design is ideal for small teams or early development stages, offering simplicity, unified database access, and easy maintenance of ACID properties. However, as the system grows, microservices provide better scalability and fault isolation. For instance, the menu service benefits from caching and read replicas due to frequent reads, while the order service requires transactional integrity and independent scaling. Additionally, microservices allow technology diversity—such as using NoSQL for the menu and PostgreSQL for orders—enabling flexibility and performance optimization as system demands evolve.

### Database-per-Service Pattern: Benefits and Complications
The database-per-service pattern in the Student Cafe system enables each microservice to own and optimize its data independently, allowing flexible schema changes and targeted performance tuning. However, this isolation introduces challenges in maintaining data consistency and handling cross-service operations. Unlike the monolithic model, where user and menu validations occurred within a single transaction, the microservices setup requires inter-service communication and orchestration, increasing complexity and potential failure points. Additionally, tasks like generating analytical reports across services become more difficult, often requiring data duplication or eventual consistency mechanisms to balance autonomy with functionality.

### When NOT to Split a Monolith
Avoiding microservices is often the wiser choice when teams are small or system domains are tightly coupled. For the Student Cafe system, a monolithic architecture remains more efficient if the development team lacks the resources to manage distributed systems, as microservices introduce added complexity through service discovery, communication, and debugging challenges. When business logic—such as real-time inventory updates affecting orders—is highly interdependent, separating services can reduce performance and increase coordination overhead. Moreover, if domain boundaries are still evolving or all functionality is handled by a single team, maintaining a monolith aligns better with Conway’s Law, ensuring simplicity, cohesion, and faster development cycles.

### Inter-Service Communication and Validation Patterns
In the Student Cafe system, inter-service communication is essential for maintaining data integrity across microservices. The order service validates user and menu data via HTTP requests to the respective services, preserving autonomy and ensuring accurate order creation. However, this introduces latency and potential network failures, as each order involves multiple service calls. To mitigate issues, the system adopts a fail-fast strategy—terminating order creation immediately if validation fails. While eventual consistency could improve performance by deferring validations, the synchronous approach is preferred here, as it prioritizes accuracy and reliability, which are critical for a real-time ordering experience.

### Resilience and Failure Handling
The menu service’s availability during order creation underscores the challenge of handling partial failures in microservices. In the current setup, if the menu service goes down, the entire order process fails—prioritizing consistency over availability. To enhance resilience, strategies like implementing circuit breakers can help the system fail gracefully, while caching menu data in the order service can allow temporary operations during downtime. Alternatively, accepting orders with “unknown” menu items for later validation promotes availability but adds complexity and risks temporary inconsistency. Each solution reflects a trade-off between reliability, consistency, and system complexity.

### Performance Optimization Through Caching
Caching offers significant performance gains for the Student Cafe system, particularly for the read-intensive menu service. By using Redis or API gateway-level response caching, frequently accessed menu data can be served from memory, reducing database load and improving response times. However, caching also introduces data consistency challenges, as updates to menu items or prices require precise invalidation strategies to prevent stale data. Similarly, caching validated user or menu data in the order service can reduce redundant network calls, but maintaining accuracy demands careful synchronization and invalidation mechanisms to balance speed with reliability.

## Conclusion
This practical demonstrated clear identification of architectural characteristics and trade-offs, effectively applying domain-driven design principles to establish logical service boundaries around user management, menu operations, and order processing. The systematic extraction process maintained full functionality while transitioning from a single codebase to independent services, implemented robust service discovery patterns using Consul for dynamic service registration and health monitoring, and successfully orchestrated the complete multi-service ecosystem using Docker Compose with proper networking and dependency management. Through hands-on implementation using Go, Docker, PostgreSQL, and Consul, this exercise provided invaluable insights into real-world distributed systems challenges—from inter-service communication and data consistency to infrastructure management and cross-service validation. This practical serves as an essential foundation for understanding enterprise-level system design and the practical realities of modern distributed software architecture.
