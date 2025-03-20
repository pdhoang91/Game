FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o oden-server ./cmd/api

# Create a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/oden-server .

# Copy configuration
COPY --from=builder /app/internal/config/config.json ./internal/config/config.json

# Expose port
EXPOSE 8080

# Command to run
CMD ["./oden-server"] 