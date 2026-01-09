# ============================================
# Stage 1: Frontend Builder
# ============================================
FROM node:20-alpine AS frontend-builder

WORKDIR /build

# Copy package files first for better caching
COPY web/package.json web/package-lock.json* ./

# Install dependencies
RUN npm ci --prefer-offline --no-audit || npm install

# Copy frontend source
COPY web/ ./

# Build frontend
RUN npm run build

# ============================================
# Stage 2: Backend Builder
# ============================================
FROM golang:1.23-alpine AS backend-builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with optimizations
# CGO_ENABLED=0 for pure Go SQLite (modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o bitwarden-backup \
    ./cmd/server

# ============================================
# Stage 3: Runtime
# ============================================
FROM node:20-alpine AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Install Bitwarden
ARG BW_CLI_VERSION=latest
RUN npm install -g @bitwarden/cli@${BW_CLI_VERSION} && \
    npm cache clean --force

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/data /app/backups /app/.tmp && \
    chown -R appuser:appgroup /app

# Copy binary from backend builder
COPY --from=backend-builder /build/bitwarden-backup ./

# Copy frontend dist from frontend builder
COPY --from=frontend-builder /build/dist ./web/dist

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Environment variables
ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/bitwarden-backup.db
ENV APP_ENV=production

# Expose port
EXPOSE 8080

# Health check (optional, uses root endpoint)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/servers || exit 1

# Run application
CMD ["./bitwarden-backup"]
