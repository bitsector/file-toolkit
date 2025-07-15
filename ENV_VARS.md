# Environment Variables Configuration

This application follows containerization best practices for configuration management:

## Configuration Priority (Highest to Lowest)

1. **OS Environment Variables** - Set by container runtime, Kubernetes, Docker, etc.
2. **`.env` file** - Loaded if present in the application directory
3. **Default Values** - Hardcoded fallbacks in the application

## Available Environment Variables

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | HTTP server port |
| `BUFFER_SIZE` | `10485760` | Multipart form buffer size in bytes (10MB) |

### Worker Pool Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `NUM_WORKERS` | `5` | Number of worker goroutines for processing conversions |
| `JOB_TIMEOUT` | `30s` | Maximum time allowed for a conversion job |
| `WORKER_RESULT_TIMEOUT` | `1s` | Timeout for sending results back to client |
| `JOB_QUEUE_TIMEOUT` | `100ms` | Timeout for submitting job to worker queue |
| `JOB_QUEUE_SIZE` | `100` | Buffer size for the job queue channel |

### File System Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `UPLOAD_PATH` | `./uploads` | Directory for temporary file uploads |

## Duration Format

Duration values support Go's duration format:
- `30s` - 30 seconds
- `1m` - 1 minute
- `100ms` - 100 milliseconds
- `1h30m` - 1 hour 30 minutes

## Usage Examples

### Local Development with .env file

```bash
# Copy example file
cp .env.example .env

# Edit values as needed
nano .env

# Run application
go run ./cmd/converter
```

### Docker Container

```bash
# Build image
docker build -t webp-converter .

# Run with environment variables
docker run -p 3000:3000 \
  -e NUM_WORKERS=10 \
  -e JOB_TIMEOUT=60s \
  -e BUFFER_SIZE=20971520 \
  webp-converter
```

### Docker Compose

```bash
# Using environment section in docker-compose.yml
docker-compose up

# Or using .env file
docker-compose --env-file .env up
```

### Kubernetes Deployment

```bash
# Apply the deployment
kubectl apply -f k8s-deployment.yaml

# Or use ConfigMap/Secrets
kubectl create configmap webp-converter-config --from-env-file=.env
```

## Performance Tuning

### For High Concurrency
- Increase `NUM_WORKERS` (e.g., 10-20)
- Increase `JOB_QUEUE_SIZE` for higher throughput (e.g., 200-500)
- Increase `BUFFER_SIZE` for larger files
- Increase `JOB_TIMEOUT` for complex conversions

### For Memory Optimization
- Decrease `NUM_WORKERS`
- Decrease `BUFFER_SIZE`
- Decrease timeouts to fail fast

### For Kubernetes
- Use ConfigMaps for non-sensitive configuration
- Use Secrets for sensitive data
- Set appropriate resource limits and requests
