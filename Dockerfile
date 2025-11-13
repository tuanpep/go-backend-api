# syntax=docker/dockerfile:1.4

# Build arguments for versioning and metadata
ARG BUILD_VERSION=unknown
ARG BUILD_DATE
ARG VCS_REF
ARG GO_VERSION=1.24

# Build stage
FROM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum* ./

# Download dependencies with BuildKit cache mount for faster rebuilds
# GOTOOLCHAIN=auto allows automatic toolchain upgrades for dependencies requiring newer Go versions
ENV GOTOOLCHAIN=auto
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy source code (needed for go mod tidy to work properly)
COPY . .

# Generate/update go.sum based on actual imports
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy && go mod verify

# Build the application with optimizations
# -ldflags='-w -s' strips debug info and symbol table
# -extldflags "-static" creates a fully static binary
# -trimpath removes file system paths from the resulting executable
ARG BUILD_VERSION
ARG BUILD_DATE
ARG VCS_REF
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags='-w -s -extldflags "-static" -X main.version=${BUILD_VERSION} -X main.buildDate=${BUILD_DATE} -X main.vcsRef=${VCS_REF}' \
    -a -installsuffix cgo \
    -o main \
    cmd/main.go

# Final stage - minimal alpine image for security and size
# Using specific version for reproducibility
FROM alpine:3.19

# Install only runtime dependencies needed
# ca-certificates for HTTPS, tzdata for timezone support
# busybox-extras provides wget for health checks (smaller than full wget)
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    busybox-extras \
    && rm -rf /var/cache/apk/* /tmp/*

# Create non-root user with specific UID/GID for security
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup -s /bin/sh appuser

# Create app directory
WORKDIR /app

# Copy the binary from builder stage with correct ownership
COPY --from=builder --chown=appuser:appgroup /app/main /app/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check using wget (from busybox-extras)
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Add labels for metadata and security scanning
LABEL org.opencontainers.image.title="go-backend-api" \
      org.opencontainers.image.description="Go Backend API - REST API built with Go" \
      org.opencontainers.image.version="${BUILD_VERSION}" \
      org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.revision="${VCS_REF}" \
      org.opencontainers.image.source="https://github.com/your-org/go-backend-api" \
      maintainer="your-email@example.com"

# Run the application
CMD ["/app/main"]
