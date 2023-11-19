# Use an official Golang runtime as a parent image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application inside the container
RUN go build -o app

# Expose a port for the application to listen on
EXPOSE 8000

# Define the command to run your application
CMD ["./app"]
