# Bitwarden Backup

自动化 Bitwarden 密码库备份和迁移工具，支持多目标备份（本地、WebDAV、S3、目标服务器）。

[![GitHub Release](https://img.shields.io/github/v/release/mingzaily/bitwarden-backup?include_prereleases)](https://github.com/mingzaily/bitwarden-backup/releases)
[![Docker Image](https://ghcr-badge.egpl.dev/mingzaily/bitwarden-backup/latest_tag?trim=major&label=Docker%20Image)](https://github.com/mingzaily/bitwarden-backup/pkgs/container/bitwarden-backup)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub Stars](https://img.shields.io/github/stars/mingzaily/bitwarden-backup?style=social)](https://github.com/mingzaily/bitwarden-backup)

![新版控制台总览截图](main.png)

## 为什么选择这个工具？

GitHub 上大多数 Bitwarden 备份方案采用 **数据库文件备份 + rclone** 的方式，直接复制 Vaultwarden/Bitwarden 的 SQLite 数据库文件。这种方式存在一些局限：

- 需要直接访问服务器文件系统
- 备份的是加密的数据库文件，恢复时依赖原有密钥
- 无法跨服务器迁移数据

**本工具采用不同的思路** —— 基于 **Bitwarden CLI 导出 JSON** 的方式：

| 特性 | 数据库备份 + rclone | 本工具（JSON 导出） |
|------|---------------------|---------------------|
| 备份方式 | 复制加密的 db 文件 | 通过 API 导出明文/加密 JSON |
| 服务器访问 | 需要文件系统权限 | 仅需 API 凭证 |
| 跨服务器迁移 | 不支持 | 原生支持 |
| 官方服务器 | 不适用 | 完全支持 |
| 备份可读性 | 不可读 | JSON 格式，可检查 |
| 恢复方式 | 替换数据库文件 | 通过 Bitwarden 导入 |

**适用场景**：
- 使用官方 Bitwarden 服务，无法访问服务器文件
- 需要在多个 Bitwarden 实例间迁移数据
- 希望备份文件可读、可验证
- 需要 Web 界面管理备份任务

## 功能特性

- 🕐 定时自动备份（支持 5/6 位 Cron 表达式）
- 📦 多种存储目标支持：本地存储、WebDAV、S3 兼容存储、目标服务器迁移
- 🖥️ Web 管理界面（简约科技风 UI）
- 🧭 工作区总览：集中查看资源、任务和最近运行状态
- 🔄 支持多个 Bitwarden 源站配置
- 📊 备份历史和日志查看（支持分页）
- 🧹 自动清理：临时文件清理 + 备份保留策略
- 🔐 AES-256-GCM 加密保护敏感凭证
- 🐳 多架构 Docker 镜像（amd64/arm64）

Web 控制台按使用流程组织：`总览`负责状态概览，`备份任务`负责编排执行，`备份资源`集中管理 `Bitwarden 源站` 与 `存储目标`，`运行记录`负责追踪结果和产物。原有 `/servers`、`/destinations`、`/tasks`、`/logs` 地址继续兼容。

## ⚠️ 安全提示

> **重要**：Web 管理界面和 API 默认需要管理员密码登录。请通过 `BITWARDEN_BACKUP_ADMIN_PASSWORD` 设置至少 8 位的管理员密码；生产环境仍建议使用更长的随机密码。会话使用 HttpOnly + SameSite Cookie，写操作启用 CSRF 校验。
>
> **强烈建议**：
> - 仅在**内网环境**或**本地**运行，不要直接暴露到公网
> - 如需远程访问，请通过 **HTTPS 反向代理** 或 VPN，并设置 `AUTH_COOKIE_SECURE=true`
> - 数据库中存储了 Bitwarden 凭证（已加密），请妥善保护 `data/` 目录

### 登录保护方式

本项目当前使用的是**应用层会话登录**，不是 Nginx Basic Auth：

1. 浏览器向 `POST /api/auth/login` 提交管理员密码。
2. 服务端签发随机的 `HttpOnly`、`SameSite=Strict` 会话 Cookie；密码本身不会写入数据库。
3. 查询接口要求有效会话，写操作还要求 `X-CSRF-Token`，用于防止跨站请求伪造。
4. 会话保存在当前进程内存中，有效期为 12 小时；服务重启后所有会话失效。当前方案适合单实例部署。

Nginx Basic Auth 可以作为外层防护，但通常不需要替代应用登录。对于只允许内网访问、应用端口不直接暴露且由 Nginx 独占入口的部署，Basic Auth 作为网关保护已经能挡住未授权访问；仍建议保留应用登录，因为它提供了 API/SPA 会话语义和 CSRF 保护。组合部署时请同时做到：

- Nginx 对外提供 HTTPS，并设置 `AUTH_COOKIE_SECURE=true`；
- 防火墙或容器网络阻止绕过 Nginx 直接访问应用端口；
- 在 Nginx 中配置 `auth_basic`/`auth_basic_user_file`，并按需增加限流。

最小化的 Nginx 入口示例（证书配置略）：

```nginx
location / {
    auth_basic "Bitwarden Backup";
    auth_basic_user_file /etc/nginx/.htpasswd;
    proxy_set_header X-Forwarded-Proto https;
    proxy_pass http://127.0.0.1:8080;
}
```

如果运行多个应用实例，内置会话不会在实例间共享；应使用单实例，或在未来接入共享会话存储并配置负载均衡粘性会话。

## 快速开始

### 方式一：Docker Compose（推荐）

创建 `docker-compose.yml` 文件：

```yaml
services:
  bitwarden-backup:
    image: ghcr.io/mingzaily/bitwarden-backup:latest
    container_name: bitwarden-backup
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
      - ./backups:/app/backups
    environment:
      BITWARDEN_BACKUP_MASTER_KEY: your-secret-key-here  # 建议使用随机字符串
      BITWARDEN_BACKUP_ADMIN_PASSWORD: "${BITWARDEN_BACKUP_ADMIN_PASSWORD:?Set this environment variable first}"
      TZ: Asia/Shanghai
```

启动服务：

```bash
export BITWARDEN_BACKUP_ADMIN_PASSWORD='replace-with-a-long-random-password'
docker compose up -d
```

访问 `http://localhost:8080` 进入管理界面。

### 方式二：Docker Run

```bash
docker run -d \
  --name bitwarden-backup \
  -p 8080:8080 \
  -e BITWARDEN_BACKUP_MASTER_KEY=your-secret-key-here \
  -e BITWARDEN_BACKUP_ADMIN_PASSWORD=replace-with-a-long-random-password \
  -e TZ=Asia/Shanghai \
  -v ./data:/app/data \
  -v ./backups:/app/backups \
  ghcr.io/mingzaily/bitwarden-backup:latest
```

### 方式三：从源码构建

```bash
git clone https://github.com/mingzaily/bitwarden-backup.git
cd bitwarden-backup

# 构建前端
cd web && npm install && npm run build && cd ..

# 构建后端
go build -o bitwarden-backup ./cmd/server

# 运行
./bitwarden-backup
```

**前置要求**：
- Go 1.25.13+
- Node.js 22+
- [Bitwarden CLI](https://bitwarden.com/help/cli/)（需全局安装：`npm install -g @bitwarden/cli`）

## 配置说明

### 环境变量

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `BITWARDEN_BACKUP_MASTER_KEY` | 否 | 加密主密钥 | 自动生成 |
| `BITWARDEN_BACKUP_ADMIN_PASSWORD` | 是 | Web 管理员密码，至少 8 位（生产环境建议更长） | 无默认值 |
| `SERVER_PORT` | 否 | 服务端口 | `8080` |
| `DB_PATH` | 否 | 数据库路径 | `/app/data/bitwarden-backup.db` |
| `AUTH_COOKIE_SECURE` | 否 | HTTPS 反向代理下启用 Secure Cookie | `false` |
| `TZ` | 否 | 时区 | `Asia/Shanghai` |

### 加密密钥管理

**推荐**：通过环境变量 `BITWARDEN_BACKUP_MASTER_KEY` 配置密钥。

如未配置，系统会自动生成并保存到 `data/.env`。

> **提示**：建议定期备份 `data/` 目录，其中包含数据库和密钥文件。

## 使用指南

### 1. 配置 Bitwarden 源站

在“备份资源 → Bitwarden 源站”页面添加 Bitwarden 服务器配置：

- **服务器 URL**: Bitwarden 服务器地址（如 `https://vault.bitwarden.com`）
- **Client ID / Client Secret**: API 密钥（在 Bitwarden 设置中获取）
- **Master Password**: 主密码（用于解锁密码库导出数据）

### 2. 配置存储目标

在“备份资源 → 存储目标”页面配置备份文件落点，支持四种类型：

- **本地存储**: 备份到 `/app/backups` 目录
- **WebDAV**: 备份到 WebDAV 服务器（如 Nextcloud）
- **S3**: 备份到 S3 兼容存储（AWS S3、MinIO、阿里云 OSS 等）
- **目标服务器**: 导入到另一个 Bitwarden 服务器

### 3. 创建备份任务

在“备份任务”页面把一个 Bitwarden 源站和一个或多个存储目标组合起来，支持 6 位 Cron 表达式（秒 分 时 日 月 周）：

```
0 0 2 * * *    # 每天凌晨 2 点
0 0 */6 * * *  # 每 6 小时
0 30 1 * * 1   # 每周一凌晨 1:30
```

### 4. 查看运行记录

“总览”页面会展示资源数量、任务状态和最近活动；“运行记录”页面可按任务筛选，并查看备份文件和详细执行过程。

## 目录结构

```
/app/
├── data/       # 数据库文件和密钥
└── backups/    # 本地备份文件
```

导出过程使用系统随机临时目录保存中间文件，并在任务结束后清理。

## 技术栈

- **后端**: Go 1.25.13, Gin, GORM, SQLite
- **前端**: Vue 3, Vite, Tailwind CSS
- **调度**: robfig/cron（支持秒级调度）
- **加密**: AES-256-GCM + PBKDF2

## 开发

```bash
# 开发模式（前端热重载 + 后端）
cd web && npm run dev &
go run ./cmd/server

# 或使用开发脚本
./dev.sh
```

### 启动故障排查

如果看到 `toolchain@v0.0.1-go1.25.13... checksum database disabled by GOSUMDB=off`，说明本机 Go 需要自动下载 Go 1.25.13，但全局配置关闭了 checksum 校验。不要将项目降级到 Go 1.25.12，直接用官方 checksum 服务启动：

```bash
GOSUMDB=sum.golang.org \
GOPROXY=https://proxy.golang.org,direct \
GOTOOLCHAIN=auto \
./dev.sh
```

如果只启动后端，将最后一行改为 `go run ./cmd/server` 即可。

如果希望永久修复当前用户的 Go 配置：

```bash
go env -w GOSUMDB=sum.golang.org
go env -w GOPROXY=https://proxy.golang.org,direct
go env -w GOTOOLCHAIN=auto
```

## License

MIT
