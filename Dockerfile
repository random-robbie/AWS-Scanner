# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aws-scanner main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 scanner && \
    adduser -D -u 1000 -G scanner scanner

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/aws-scanner .

# Create output directory
RUN mkdir -p /app/output && \
    chown -R scanner:scanner /app

# Switch to non-root user
USER scanner

# Set entrypoint
ENTRYPOINT ["./aws-scanner"]

# Default command (can be overridden)
CMD ["--list", "list.txt"]
