# Stage 1: Build the application
FROM golang:1.22-alpine3.19 AS builder

# Set up the working directory
WORKDIR /app

# Copy the Go application code
#COPY main.go .
COPY * .

# Build the application binary
RUN go build -o auth-service main.go

# Stage 2: Create the final image with just the binary
FROM alpine:3.16

# Install necessary certificates
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/auth-service /auth-service

# Expose the port on which the service will run
EXPOSE 8081

# Run the application
ENTRYPOINT ["/auth-service"]

