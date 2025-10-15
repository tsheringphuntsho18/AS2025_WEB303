# Practical 5: Microservices Architecture - Student Cafe System

## Overview

This project demonstrates the transformation of a monolithic Student Cafe application into a distributed microservices architecture using Go, Docker, PostgreSQL, and Consul for service discovery.

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

    Inter-Service Communication:
    Order Service → User Service (Validate user exists)
    Order Service → Menu Service (Validate menu items & get prices)
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

### Service Health Endpoints

- User Service: http://localhost:8081/health
- Menu Service: http://localhost:8082/health
- Order Service: http://localhost:8083/health

## Architecture Benefits

### Scalability

- Each service can be scaled independently
- Database load is distributed across multiple instances
- Services can use different resource allocations

### Maintainability

- Clear service boundaries reduce complexity
- Independent deployment cycles
- Technology diversity possible per service

### Resilience

- Failure in one service doesn't bring down entire system
- Circuit breaker patterns can be implemented
- Health checks enable automatic recovery

### Team Organization

- Different teams can own different services
- Parallel development possible
- Clear API contracts between services

## Reflection Essay

### Architectural Decision Analysis: Monolith vs Microservices in Practice

The transformation of the Student Cafe system from a monolithic architecture to microservices provides valuable insights into the practical implications of distributed system design. This reflection examines the key architectural decisions, trade-offs, and lessons learned from implementing a microservices approach for this specific use case.

#### Monolith vs Microservices: Context-Driven Decision Making

For the Student Cafe system, the choice between monolith and microservices depends heavily on organizational and technical context. The monolithic approach offers significant advantages for small teams or early-stage applications. In the original monolith, all functionality—user management, menu operations, and order processing—resided in a single codebase with shared database access. This simplicity enables rapid development, straightforward debugging, and easier transaction management. ACID properties are maintained naturally when all operations occur within the same database transaction scope.

However, the microservices approach demonstrates its value when considering scalability patterns specific to a cafe system. Menu items are read frequently but updated infrequently, making the menu service an ideal candidate for heavy caching and read replicas. Conversely, order processing involves complex validation logic and benefits from isolation to prevent failures in one area from affecting others. User management sits between these extremes, requiring moderate scaling with strong consistency for authentication purposes.

The microservices architecture also enables technology diversity—the menu service could potentially use a NoSQL database for better read performance, while the order service maintains PostgreSQL for transactional integrity. This flexibility becomes crucial as system requirements evolve and different components face distinct performance challenges.

#### Database-per-Service Pattern: Benefits and Complications

The database-per-service pattern implemented in this project illustrates both the power and complexity of data isolation in microservices. Each service maintains complete ownership of its data schema, enabling independent evolution and optimization. The user service can modify its user table structure without coordinating with other teams, while the menu service can implement category-specific indexing strategies.

However, this separation introduces significant challenges around data consistency and cross-service queries. In the monolithic version, creating an order with user and menu validation occurred within a single transaction. The microservices version requires careful orchestration: the order service must validate user existence via HTTP calls to the user service, then validate each menu item through the menu service, before finally creating the order record. This process lacks the atomic guarantees of database transactions and introduces potential failure points.

The trade-off becomes apparent in scenarios requiring complex reporting or analytics. Generating a report showing "orders by user demographics and menu category preferences" would be straightforward in the monolith with JOIN queries but requires either service-to-service communication or eventual consistency patterns in the microservices approach. Data duplication—such as storing user email snapshots in orders—becomes a necessary evil to maintain service independence while providing essential functionality.

#### When NOT to Split a Monolith

The decision to avoid microservices should be driven by pragmatic considerations rather than architectural preferences. Small teams (fewer than 8-10 developers) often lack the operational expertise to manage distributed systems effectively. The overhead of service discovery, inter-service communication, and distributed debugging can overwhelm development velocity.

Additionally, tightly coupled business logic presents strong arguments against decomposition. If the cafe system included complex inventory management where menu availability directly influenced order processing in real-time, the communication overhead between services might outweigh the benefits of separation. Similarly, systems with unclear domain boundaries should remain monolithic until usage patterns and team understanding mature sufficiently to identify proper service boundaries.

The "Conway's Law" principle also applies—organizations should not create service boundaries that don't align with team structures. If a single team maintains all cafe functionality, creating artificial service boundaries adds complexity without organizational benefit.

#### Inter-Service Communication and Validation Patterns

The order service's validation of user existence demonstrates a fundamental microservices pattern: service-to-service communication for data validation. Rather than accessing user data directly, the order service makes HTTP GET requests to `http://user-service:8081/users/{id}`. This approach maintains service autonomy while ensuring data integrity.

However, this pattern introduces latency and potential failure modes. Each order creation requires at least two additional network calls (user validation and menu item validation), increasing response time and failure probability. The system handles this by failing fast—if user validation fails, the order creation aborts immediately rather than proceeding with invalid data.

A more sophisticated approach might implement eventual consistency patterns, accepting orders with unvalidated users and rectifying inconsistencies through background processes. However, for a cafe system where order accuracy is crucial, the synchronous validation approach provides better user experience despite performance costs.

#### Resilience and Failure Handling

The question of menu service availability during order creation highlights a critical microservices challenge: partial system failures. In the current implementation, if the menu service becomes unavailable, order creation fails completely. This represents a trade-off between data consistency and system availability.

Several strategies could improve resilience: implementing circuit breaker patterns to fail gracefully when services are unavailable, maintaining cached menu data in the order service for emergency scenarios, or designing the system to accept orders with "unknown" menu items for later validation. Each approach involves different consistency and complexity trade-offs.

#### Performance Optimization Through Caching

Caching presents excellent opportunities for performance improvement in this architecture. The menu service exhibits classic read-heavy patterns ideal for Redis implementation. Menu items could be cached with appropriate TTL values, dramatically reducing database load and improving response times. The API gateway could implement response caching for menu endpoints, serving frequently requested items directly from memory.

However, caching introduces consistency challenges—cached menu prices must be invalidated when items are updated, requiring careful cache invalidation strategies. For the order service, caching validated user and menu data could reduce redundant validation calls, but requires sophisticated invalidation logic to maintain data accuracy.

#### Conclusion

The Student Cafe microservices implementation demonstrates that architectural decisions must be driven by specific organizational and technical requirements rather than industry trends. While the microservices approach provides clear benefits in scalability, team autonomy, and technology diversity, it introduces significant complexity in inter-service communication, data consistency, and operational overhead. The success of such architectures depends heavily on team maturity, organizational structure, and the specific characteristics of the business domain being modeled.


This practical demonstrates the fundamental principles of microservices architecture while highlighting both the benefits and challenges of distributed systems design.
