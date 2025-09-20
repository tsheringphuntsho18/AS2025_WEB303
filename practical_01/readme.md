# Practical_01: From Foundational Setup to Inter-Service Communication

## Overview

This practical demonstrates a simple microservices architecture using Go and gRPC. It consists of two services:

- **Greeter Service**: Responds to greeting requests.
- **Time Service**: Provides the current server time.

Both services communicate using Protocol Buffers (`.proto` files) and are containerized using Docker. The practical uses `docker-compose` for orchestration.

---

## Directory Structure

```
practical_01/
│
├── docker-compose.yml
├── go.mod
├── go.sum
├── readme.md
│
├── greeter-service/
│   ├── Dockerfile
│   └── main.go
│
├── time-service/
│   ├── Dockerfile
│   └── main.go
│
└── proto/
    ├── greeter.proto
    ├── time.proto
    └── gen/
        ├── greeter.pb.go
        ├── greeter_grpc.pb.go
        ├── time.pb.go
        └── time_grpc.pb.go
```

---

## Approach

### Part 1: Foundational Development Environment Setup

1. **Go Installation**
   - Downloaded and installed Go from [https://go.dev/dl/](https://go.dev/dl/).
   - Verified installation with `go version` and `go env`.

![Go version](/assets/go.png)

2. **Protocol Buffers & gRPC Tools**
   - Installed the Protobuf compiler (`protoc`) from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases).
   - Installed Go plugins for Protobuf and gRPC:
     ```sh
     go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
     go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
     export PATH="$PATH:$(go env GOPATH)/bin"
     ```
   - Verified installation by running `protoc --version`.

![protoc verification](/assets/protoc.png)

3. **Docker Installation**
   - Installed Docker Desktop from [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/).
   - Verified Docker with `docker run hello-world`.

![docker verification](/assets/dockerhello.png)

---

### Part 2: Building and Orchestrating Communicating Microservices

1. **Project Structure & Service Contracts**
   - Created the required directory structure.
   - Defined `greeter.proto` and `time.proto` in the `proto/` directory, specifying service interfaces and messages.

2. **Generating Go Code from Protobuf**
   - Ran the following command to generate Go code:
     ```sh
     protoc --go_out=./proto/gen --go_opt=paths=source_relative \
            --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
            proto/*.proto
     ```
   - Verified that `proto/gen/` contained the four generated `.go` files.

3. **Implementing the Microservices**
   - **time-service:** Implemented a gRPC server that returns the current time.
   - **greeter-service:** Implemented a gRPC server that calls the time-service and returns a greeting with the current time.

4. **Containerization**
   - Wrote Dockerfiles for both services using multi-stage builds for minimal images.
   - Ensured correct `COPY` paths for `go.mod`, `go.sum`, and source files(the four go files).

5. **Orchestration with Docker Compose**
   - Created `docker-compose.yml` to build and run both services.
   - Set up service dependencies so that greeter-service waits for time-service.

---

## Run and Verify

**1. Run Docker Compose:** From the root of your practical_01 directory, run:  
```
sudo docker-compose up --build
```

**Expected Outcome:** 

![docker build](/assets/build.png)

**2. Test the Endpoint:** To test the flow, we'll use grpcurl.  
```bash
grpcurl -plaintext \
    -import-path ./proto -proto greeter.proto \
    -d '{"name": "WEB303 Student"}' \
    0.0.0.0:50051 greeter.GreeterService/SayHello
```


**Expected Outcome:** 

![grpc call](/assets/grpcurl.png)

## Steps to Run

1. **Install Dependencies**
   ```sh
   go mod tidy
   ```

2. **Install Protobuf Plugins**
   ```sh
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

3. **Generate Go Code from Protobuf**
   ```sh
   cd practical_01
   protoc --go_out=./proto/gen --go_opt=paths=source_relative \
          --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
          proto/*.proto
   ```

4. **Build and Run with Docker Compose**
   ```sh
   sudo docker-compose up --build
   ```

5. **Test the Endpoint**
   - Install `grpcurl` if not already installed.
   - In a new terminal, run:
     ```sh
     grpcurl -plaintext \
         -import-path ./proto -proto greeter.proto \
         -d '{"name": "WEB303 Student"}' \
         0.0.0.0:50051 greeter.GreeterService/SayHello
     ```
   - **Expected Output:**
     ```json
     {
       "message": "Hello WEB303 Student! The current time is 2025-07-24T09:45:00Z"
     }
     ```

---

## Challenges Encountered
I encountered the major challenges during `docker-compose up --build`

**1. Permission Denied**  
Docker need root user permission to be run. So run:

```bash 
sudo docker-compose up --build
```   
    
**2. Folder Renaming Issues**  
Renaming the project folder from `practical-one` to `practical_01` caused runtime panics in the services. The generated protobuf files still contained the old module paths, resulting in service crashes immediately after startup.

I have learned that after directory and module renaming, the protobuf files had to be regenerated to align import paths with the Go module. Otherwise, services crashed due to protobuf initialization errors. So I deleted the old four go code and regenerated Go Code.  
```bash
protoc --go_out=./proto/gen --go_opt=paths=source_relative \
    --go-grpc_out=./proto/gen --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

**3. Proto Code Generation Paths**  
When we generate the Go code, it goes into the `./proto/gen/proto` directory and gave error when running the docker build command. So I ensured the generated code was placed in the correct directory (`./proto/gen/`).

**4. Go Module Import Errors**  
This occurred because the Go module name in `go.mod` did not match the import paths used in the code or the folder structure.  
Since my go.mod says `module practical_01`, Go expects the module root to be the directory where go.mod should be. But in my case, it was inside time-service and greeter-service directory.  
So I moved the go.mod and go.sum to the root directory that is practical_01 and i deleted the another go.mod and go.sum that was inside greeter-service.

**5. Dockerfile**  
Go cannot see your `proto/gen` folder inside Docker because Dockerfile `COPY` commands failed as incorrect path were used. So with `COPY . .` it copied the whole project root, so contains `proto/gen` exists before Go build.

**6.Docker Volume / Container Configuration Errors**   
Errors like `KeyError: 'ContainerConfig'` appeared, likely caused by leftover containers or corrupted Docker metadata. It required stopping all containers, pruning unused volumes/images, and rebuilding. 

```bash 
sudo docker system prune -af  
sudo docker volume prune -f
``` 
This removes all stopped containers, dangling images, and unused volumes.  
Rebuild everything clean  

```bash
sudo docker-compose up --build --force-recreate
```

---

## Notes
- The main challenges were module path mismatches.
- The Directory structure should be well maintain. 
- Regenerate the Go files after any changes to `.proto` files.
- Both services are designed to be simple and demonstrate gRPC communication and Docker-based deployment.
- This practical lays the foundation for deploying microservices to Kubernetes in future exercises.

---
