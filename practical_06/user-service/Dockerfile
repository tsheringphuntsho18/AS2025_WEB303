FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod + go.sum
COPY user-service/go.mod user-service/go.sum ./

# Copy protos
COPY student-cafe-protos ./student-cafe-protos

# Download modules
RUN go mod download

# Copy actual service source code
COPY user-service/. .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service

# Final stage
FROM alpine:latest
COPY --from=builder /user-service /user-service
ENTRYPOINT ["/user-service"]
