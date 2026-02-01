# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /src

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o ./bin/server ./cmd/api

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /src/bin/server /app/server

# Copy the docs directory for static file serving
COPY --from=builder /src/docs /app/docs

# Expose port (Zeabur will set PORT env var)
EXPOSE 8080

# Run the application
CMD ["/app/server"]
