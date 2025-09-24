# Practical_02 Report: API Gateway with Service Discovery

## Overview
This practical demonstrates the implementation of an API Gateway with Service Discovery using Consul for service registration and health checking.

## Project Structure
```
practical_02/
└── readme.md
```

## Repository
**Full working code repository:** [https://github.com/tsheringphuntsho18/go-microservices-demo](https://github.com/tsheringphuntsho18/go-microservices-demo)

## Screenshots

### 1. Consul UI - Service Registration and Health Status
![Consul UI Services](/assets/consul_ui.png)
*Screenshot of the Consul UI showing both services registered and healthy.*

### 2. API Requests via Postmam

Test the users service
![Postman Requests](/assets/postman1.png)

Test the products service
![Postman Requests](/assets/postman2.png)

*Screenshot demonstrating successful API calls to both services through the gateway*

### 3. API Gateway Terminal Output
![API Gateway Logs](/assets/api_gateway1.png)

![API Gateway Logs](/assets/api_gateway2.png)
*Terminal output showing successful request routing and service communication*

---

## Conclusion

This practical successfully demonstrates the implementation of a microservices architecture using Go with Consul for service discovery and an API Gateway for centralized routing. The screenshots confirm that both services are properly registered with Consul and maintaining healthy status, while the Postman tests validate successful API communication through the gateway. The terminal outputs show effective request routing and service-to-service communication, proving that the distributed system architecture is functioning correctly. This implementation provides a solid foundation for scalable microservices development with proper service discovery patterns.

---
