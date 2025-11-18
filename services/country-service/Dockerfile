# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files from root and shared
COPY go.mod go.sum ./
COPY shared/go.mod shared/go.sum ./shared/
COPY services/country-service/go.mod services/country-service/go.sum ./services/country-service/

# Download dependencies
WORKDIR /app/services/country-service
RUN go mod download

# Copy shared code
WORKDIR /app
COPY shared ./shared

# Copy service code
COPY services/country-service ./services/country-service

# Build the application
WORKDIR /app/services/country-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/country-service/main .

# Copy .env file (optional, for local dev)
COPY .env .env

# Expose port
EXPOSE 8082

# Run the binary
CMD ["./main"]
