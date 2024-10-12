# Use a minimal base image for efficiency
FROM golang:1.20-alpine AS builder

# Set the work directory
WORKDIR /app

# Copy the Go source code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Create a new image based on the builder image
FROM alpine:3.20

# Copy the built binary from the builder image
COPY --from=builder /app/main /app/main

# Create a non-root user
RUN adduser -D appuser

# Set the owner of the binary and its directory to the non-root user
USER appuser
WORKDIR /app

# Expose the port your application listens on (if applicable)
EXPOSE 8080

# Command to run the application as the non-root user
CMD ["/app/main"]
