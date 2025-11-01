# Multi-stage Dockerfile for beady
# Stage 1: Build the Go application
FROM golang:1.24.9-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o beady cmd/beady/*.go

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Install Beads CLI for write operations
RUN apk --no-cache add go git && \
    export GOPATH=/tmp/go && \
    go install github.com/steveyegge/beads@latest && \
    mv /tmp/go/bin/beads /usr/local/bin/bd && \
    rm -rf /tmp/go && \
    apk del go git

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/beady /app/beady

# Copy assets (if any)
COPY --from=builder /build/assets /app/assets

# Create data directory for Beads database
RUN mkdir -p /data/.beads

# Set environment variables
ENV PORT=8080
ENV BEADS_DIR=/data/.beads

# Expose port
EXPOSE 8080

# Run beady
CMD ["/app/beady", "--port", "8080"]
