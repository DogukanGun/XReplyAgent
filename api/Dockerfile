# Use the official Go image as build stage
FROM golang:1.23.2-alpine AS builder

# Build arguments for configuration
ARG SWAGGER_HOST=api.nexarb.com
ARG ENV=prod

# Set working directory
WORKDIR /app

# Copy go mod files from parent directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Install swag tool for generating swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs with dynamic host replacement
RUN swag init -g ./api/main.go -o ./api/docs
RUN sed -i "s/localhost:3002/${SWAGGER_HOST}/g" ./api/docs/swagger.json

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./api

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create api directory and copy Firebase service account file
RUN mkdir -p ./api
COPY --from=builder /app/api/firebase-service-account.json ./api/

# Change ownership of files to appuser
RUN chown -R appuser:appuser /root/

USER appuser

# Expose port 3002
EXPOSE 3002

# Set environment variables
ENV GIN_MODE=release
ENV MONGODB_URI=mongodb://localhost:27017
ENV ENV=prod

# Run the binary
CMD ["./main"]