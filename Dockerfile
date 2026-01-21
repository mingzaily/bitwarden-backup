# ============================================
# Stage 1: Frontend Builder
# ============================================
FROM --platform=$BUILDPLATFORM node:20-alpine AS frontend-builder

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
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS backend-builder

# 接收目标平台参数
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with optimizations
# 使用 Go 原生交叉编译，而非 QEMU 模拟
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
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

# 使用 node:alpine 内置的 node 用户 (UID 1000)
# 无需创建新用户，直接复用

WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/data /app/backups /app/.tmp && \
    chown -R node:node /app

# Copy binary from backend builder
COPY --from=backend-builder /build/bitwarden-backup ./

# Copy frontend dist from frontend builder
COPY --from=frontend-builder /build/dist ./web/dist

# Set ownership
RUN chown -R node:node /app

# Switch to non-root user
USER node

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
