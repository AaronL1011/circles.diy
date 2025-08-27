# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for Go modules that might need it
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata wget && \
    rm -rf /var/cache/apk/*

# Create a non-root user with no shell and no home directory
RUN adduser -D -s /sbin/nologin -H appuser

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files
COPY --from=builder /app/index.html .

# Create directory for feedback file with proper permissions
RUN mkdir -p /app/data && \
    chown appuser:appuser /app/data && \
    chmod 750 /app/data

# Set proper permissions for the application binary
RUN chmod 755 /app/main && \
    chmod 644 /app/index.html

# Switch to non-root user
USER appuser

# Set environment variables
ENV PORT=8080
ENV GIN_MODE=release
ENV TZ=Australia/Sydney

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Run the application with additional security
CMD ["./main"]