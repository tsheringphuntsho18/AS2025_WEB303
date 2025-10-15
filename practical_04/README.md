# Practical_04 Report: Kubernetes Microservices with Kong Gateway & Resilience Patterns

## Overview

The Student Cafe is a microservices-based web application deployed on Kubernetes that demonstrates modern DevOps including service discovery, advanced API gateway management with Kong, and containerized deployment.

## Architecture Components

### 1. Frontend Service (React UI)

- **Technology**: React.js with modern hooks
- **Port**: 80 (served via Nginx)
- **Docker Image**: `cafe-ui:v2`
- **Features**:
  - Interactive menu display
  - Shopping cart functionality
  - Order placement with real-time feedback
  - Responsive design with emojis and modern UI

### 2. Food Catalog Service (Go)

- **Technology**: Go with Chi router
- **Port**: 8080
- **Docker Image**: `food-catalog-service:latest`
- **Endpoints**:
  - `GET /health` - Health check
  - `GET /items` - Retrieve menu items
- **Data**: Coffee ($2.50), Sandwich ($5.00), Muffin ($3.25)

### 3. Order Service (Go)

- **Technology**: Go with Chi router
- **Port**: 8081
- **Docker Image**: `order-service:v2`
- **Endpoints**:
  - `GET /health` - Health check
  - `POST /orders` - Create new orders
- **Features**: Inter-service communication with catalog service

### 4. Infrastructure Services

- **Kong API Gateway**: Routes external requests to appropriate services
- **Consul**: Service discovery and health monitoring
- **Kubernetes**: Container orchestration and service management

## Deployment Architecture

```
Internet → Kong Proxy → Kubernetes Services → Pods
    ↓
Port Forward (8080) → Kong Gateway
    ↓
API Routes:
- /api/catalog/* → Food Catalog Service
- /api/orders/* → Order Service
- /* → Cafe UI Service
```

## Issues Encountered and Solutions

### Issue 1: Minikube Service Access

**Problem**: `minikube service` command failed with QEMU builtin network

![minikube serevice failed](/assets/practical4Screenshots/minikube-fail.png)


**Solution**: Used kubectl port-forwarding as alternative

```bash
kubectl port-forward service/kong-kong-proxy 8080:80 -n student-cafe
```
![access application](/assets/practical4Screenshots/access-application.png)


### Issue 2: React Build and Deployment

**Problem**: UI showing blank page due to improper React build
![blank ui](/assets/practical4Screenshots/blank.png)


**Solution**:

1. Rebuilt React application with production optimizations
2. Created new Docker image with proper static file serving
3. Updated Kubernetes deployment with new image version

### Issue 3: Service Discovery Failure

**Problem**: Order service couldn't find food-catalog-service through Consul  
Resulting the failure of place order functionality.
![Order fail](/assets/practical4Screenshots/order-fail.png)

![food-catalog-service not found](/assets/practical4Screenshots/postman-fail.png)

**Root Cause**: Consul health checks were failing, preventing service discovery

**Solution**: Implemented fallback mechanism in order service

```go
// Modified order service to use Kubernetes service name as fallback
catalogAddr, err := findService("food-catalog-service")
if err != nil {
    log.Printf("Warning: Could not find catalog service (%v), but continuing with order", err)
    catalogAddr = "http://food-catalog-service:8080" // Kubernetes DNS fallback
}
```

## Deployment Commands

### Building Services

```bash
# Build Docker images using minikube's Docker daemon
eval $(minikube docker-env)
docker build -t cafe-ui:v2 ./cafe-ui
docker build -t order-service:v2 ./order-service
docker build -t food-catalog-service:latest ./food-catalog-service
```

### Deployment

```bash
# Apply Kubernetes manifests
kubectl apply -f app-deployment.yaml
kubectl apply -f kong-ingress.yaml

# Update deployments with new images
kubectl set image deployment/cafe-ui-deployment cafe-ui=cafe-ui:v2 -n student-cafe
kubectl set image deployment/order-deployment order-service=order-service:v2 -n student-cafe
```

### Access Application

```bash
# Start port forwarding
kubectl port-forward service/kong-kong-proxy 8080:80 -n student-cafe

# Access application
# Open browser to: http://localhost:8080
```


## Screenshots

### 1. React Frontend - Food Menu Display

![Food Menu Interface](/assets/practical4Screenshots/food-menu-display.png)
_Screenshot showing the React frontend with the food menu items and shopping cart interface_

### 2. Successful Order Placement

![Order Success](/assets/practical4Screenshots/order-placement-success.png)
_Screenshot demonstrating successful order placement with confirmation message showing order ID and total amount_

### 3. Kubernetes Pods Status

![Kubernetes Pods](/assets/practical4Screenshots/kubectl-get-pods.png)
_Screenshot showing all running pods including cafe-ui, food-catalog, order-service, consul, and kong components_

### 4. Kubernetes Services Overview

![Kubernetes Services](/assets/practical4Screenshots/kubectl-get-services.png)
_Screenshot showing all services with their cluster IPs, external IPs, ports, and service types (ClusterIP, LoadBalancer, NodePort)_

### 5. API Testing with curl

![api testing](/assets/practical4Screenshots/postman-success.png)
_Screenshot showing successful API responses from both catalog and order services_

---

## Conclusion
I have successfully completed the assigned exercise to identify and fix the broken order submission functionality. The task involved:

**Task 1 - Deploy and Test Application:**
- Successfully deployed all microservices components on Kubernetes
- Accessed the React frontend through Kong gateway at http://localhost:8080
- Initially observed order submission failures as expected

**Task 2 - Debug and Identify Issues:**
- **Root Cause Identified:** Service discovery failure between order service and food-catalog service
- **Issue Details:** Consul health checks were preventing proper service registration and discovery
- **Solution Implemented:** Added Kubernetes DNS fallback mechanism when Consul discovery fails


**Final Status**: All services operational, orders can be placed successfully
