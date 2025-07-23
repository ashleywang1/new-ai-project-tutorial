# syntax=docker/dockerfile:1

##
## Build Stage
##
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

# Accept build arguments for multi-arch support
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Set build environment with dynamic architecture support
ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations for target architecture
RUN go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -installsuffix cgo \
    -o apiserver \
    ./cmd/apiserver/

# Verify the binary was built
RUN ls -la apiserver

##
## Runtime Stage
##
FROM gcr.io/distroless/static:nonroot

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /build/apiserver /apiserver

# Use non-root user
USER nonroot:nonroot

# Set runtime environment
ENV PORT=8080 \
    TZ=UTC

# Expose the application port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/apiserver"] 