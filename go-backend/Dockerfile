# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Enable Go modules
ENV GO111MODULE=on

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o /app/app

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Set ownership and permissions
RUN chown -R appuser:appgroup /app/app && chmod +x /app/app

# Switch to non-root user
USER appuser

# Expose a port for the application to listen on
EXPOSE 8000

# Define the command to run your application
CMD ["./app"]
