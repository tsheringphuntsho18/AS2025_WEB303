FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy the entire project (so proto/gen is included)
COPY . .

# Move into greeter-service to resolve dependencies
WORKDIR /app/greeter-service
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["/app/server"]