FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy the whole project (so proto/gen is available)
COPY . .

# Move into time-service for module context
WORKDIR /app/time-service

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["/app/server"]
