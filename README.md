# Bitwarden Backup

自动化 Bitwarden 密码库备份和迁移工具

## 功能特性

- ✅ 定时自动备份（支持 Cron 表达式）
- ✅ 自动迁移到官方服务器（可选）
- ✅ Web 管理界面
- ✅ 支持多服务器配置
- ✅ 备份历史和日志查看
- ✅ Docker 容器化部署

## 技术栈

- **后端**: Go 1.21 + Gin
- **数据库**: SQLite
- **调度器**: robfig/cron
- **前端**: 原生 HTML/CSS/JavaScript
- **部署**: Docker + Docker Compose

## 快速开始

### 使用 Docker Compose（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd bitwarden-backup

# 启动服务
docker-compose up -d

# 访问 Web 界面
open http://localhost:8080
```

### 本地开发

```bash
# 安装依赖
go mod download

# 运行服务
go run cmd/server/main.go
```

## 配置说明

### 环境变量

- `SERVER_PORT`: 服务端口（默认: 8080）
- `DB_PATH`: 数据库路径（默认: ./data/bitwarden-backup.db）
- `LOG_LEVEL`: 日志级别（默认: info）
