# Practical_03 Report: Full-Stack Microservices with gRPC, Databases, and Service Discovery

## Overview
This practical demonstrates a complete microservices ecosystem from the ground up with gRPC communication, database management, and service discovery.

## Project Structure
```
practical_02/
└── readme.md
```
## Repository
**Full working code repository:** [https://github.com/tsheringphuntsho18/Full-Stack-Microservices
](https://github.com/tsheringphuntsho18/Full-Stack-Microservices
)

## Architecture Implemented

The system consists of the following components:

1. **API Gateway**: HTTP entry point that translates requests to gRPC calls
2. **Service Discovery (Consul)**: Central registry for service locations
3. **Users Service**: Microservice managing user data with PostgreSQL
4. **Products Service**: Microservice managing product data with PostgreSQL
5. **Databases**: Separate PostgreSQL instances for each service


## Service Discovery Verification

Consul UI accessible at http://localhost:8500 showing:

- users-service registered and healthy
- products-service registered and healthy
- Services discoverable by API Gateway

## Screenshots

### 1. Terminal output
![docker output](/assets/practical3Screenshots/terminal.png)
 *Terminal showing successful startup of all microservices, databases, and Consul via docker-compose.*

### 2. Consul UI - Service Registration and Health Status
![Consul UI Services](/assets/practical3Screenshots/consul.png)
*Screenshot of the Consul UI showing both services registered and healthy.*

### 3. API Requests via Postmam
#### Create User Request
![Postman Requests](/assets/practical3Screenshots/createuser.png)
*POST request to create a new user with JSON response showing user details.*

#### Retrieve User Request
![Postman Requests](/assets/practical3Screenshots/retriveuser.png)
*GET request retrieving user by ID, demonstrating successful database persistence.*

#### Create Product Request
![Postman Requests](/assets/practical3Screenshots/createproduct.png)
*POST request to create a new product with JSON response showing product details.*

#### Retrieve Product Request
![Postman Requests](/assets/practical3Screenshots/retriveproduct.png)
*GET request retrieving product by ID from the products microservice.*

#### Aggregate User and Product Data
![Postman Requests](/assets/practical3Screenshots/retriveusernproduct.png)
*GET request demonstrating the composite endpoint aggregating data from both services.*

---

## Conclusion

This practical successfully demonstrates a complete microservices architecture with:

- Service-to-service gRPC communication
- Individual database per service
- Service discovery with Consul
- API Gateway as single entry point
- Data aggregation across services
- Containerized deployment with Docker Compose

The implementation showcases industry-standard microservices patterns and provides a scalable foundation
