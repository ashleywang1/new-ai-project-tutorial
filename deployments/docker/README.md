# Docker Deployment

This directory contains the Docker configuration for building and running the AI Project Tutorial API server.

## Files

- `Dockerfile` - Multi-stage Docker build optimized for Go applications
- `.dockerignore` - Excludes unnecessary files from build context
- `README.md` - This documentation

## Quick Start

### Building the Image

```bash
# From project root
make docker-build
```

### Running Locally

```bash
# Build and run in one command
make docker-run

# Or run manually
docker run --rm -p 8080:8080 ai-project-tutorial/apiserver:latest
```

### Multi-Architecture Build

The project supports building Docker images for multiple architectures (amd64/arm64) using Docker Buildx.

```bash
# Set up buildx builder (one-time setup)
make docker-buildx-setup

# Build for local platform only (loads into Docker)
make docker-build-local-platform

# Build for amd64 and arm64 (saves to ./bin/multiarch-image.tar)
make docker-build-multi

# Build and push multi-arch to registry
make docker-build-multi-push
```

**Note**: Multi-platform builds cannot be loaded directly into Docker due to manifest list limitations. The `docker-build-multi` target saves the image to an OCI tar file that can be loaded manually if needed.

## Image Details

- **Base Image**: `gcr.io/distroless/static:nonroot` (minimal, secure runtime)
- **Build Image**: `golang:1.24.5-alpine` (optimized for Go builds)
- **Size**: ~10-15MB (static binary + minimal runtime)
- **Security**: Runs as non-root user, no shell access
- **Architecture**: Supports amd64 and arm64

## Optimization Features

- **Layer Caching**: Go modules downloaded separately for better cache efficiency
- **Static Binary**: CGO disabled, fully static linking
- **Minimal Runtime**: Distroless image with only essential runtime files
- **Security**: Non-root user, no unnecessary packages

## Health Checks

The container exposes health endpoints:
- `GET /health` - Basic health status
- `GET /ready` - Readiness check for Kubernetes

For Kubernetes deployments, configure liveness and readiness probes to use these endpoints.

## Environment Variables

- `PORT`: Server port (default: 8080)
- `TZ`: Timezone setting (default: UTC)

## Cleanup

```bash
# Clean Docker resources
make docker-clean

# Remove buildx builder
make docker-buildx-clean

# Clean everything (build artifacts + Docker + buildx)
make clean-all
``` 