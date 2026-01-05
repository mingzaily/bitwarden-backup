# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=1 GOOS=linux go build -o bitwarden-backup ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装必要的依赖
RUN apk --no-cache add ca-certificates nodejs npm sqlite

# 安装 Bitwarden CLI
RUN npm install -g @bitwarden/cli

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bitwarden-backup .
COPY --from=builder /app/web ./web

# 创建数据目录
RUN mkdir -p /app/data /app/backups

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV SERVER_PORT=8080
ENV DB_PATH=/app/data/bitwarden-backup.db

# 启动应用
CMD ["./bitwarden-backup"]
