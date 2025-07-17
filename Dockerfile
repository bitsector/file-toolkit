# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o file-toolbox ./cmd/converter

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 appgroup && adduser -u 1001 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/file-toolbox .

# Create uploads directory and set permissions
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app && chmod 750 /app/uploads

# Switch to non-root user
USER appuser

# Set environment variables with defaults
ENV PORT=3000
ENV BUFFER_SIZE=10485760
ENV NUM_WORKERS=5
ENV JOB_TIMEOUT=30s
ENV WORKER_RESULT_TIMEOUT=1s
ENV JOB_QUEUE_TIMEOUT=100ms
ENV JOB_QUEUE_SIZE=100
ENV UPLOAD_PATH=./uploads

# Expose port
EXPOSE 3000

# Run the application
CMD ["./file-toolbox"]
