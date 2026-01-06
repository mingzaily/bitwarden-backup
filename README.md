# Bitwarden Backup

自动化 Bitwarden 密码库备份和迁移工具，支持多目标备份（本地、WebDAV、目标服务器）。

## 功能特性

- 定时自动备份（支持 Cron 表达式）
- 多备份目标支持：本地存储、WebDAV、目标服务器
- Web 管理界面
- 支持多源服务器配置
- 备份历史和日志查看
- 远程备份后自动清理临时文件
- **🔒 AES-256-GCM 加密保护敏感凭证**

## 前置要求

- [Bitwarden CLI](https://bitwarden.com/help/cli/) 已安装并配置
- Go 1.21+（本地运行）或 Docker（容器运行）

## 快速开始

### 方式一：Docker Compose（推荐）

```bash
git clone https://github.com/mingzaily/bitwarden-backup.git
cd bitwarden-backup
docker compose up -d
```

### 方式二：本地编译运行

```bash
git clone <repository-url>
cd bitwarden-backup
go build -o bitwarden-backup ./cmd/server
./bitwarden-backup
```

## 配置说明

### 环境变量

| 变量                           | 说明                   | 默认值                     | 必需 |
| ------------------------------ | ---------------------- | -------------------------- | ---- |
| `SERVER_PORT`                  | 服务端口               | 8080                       | 否   |
| `DB_PATH`                      | 数据库路径             | ./data/bitwarden-backup.db | 否   |
| `LOG_LEVEL`                    | 日志级别               | info                       | 否   |
| `BITWARDEN_BACKUP_MASTER_KEY`  | 加密主密钥             | 自动生成并保存到 .env      | 自动 |

### 🔒 加密密钥管理

**重要说明**：系统会自动管理加密密钥，确保数据安全和持久性。

#### 密钥加载优先级

1. **环境变量** `BITWARDEN_BACKUP_MASTER_KEY`（推荐用于生产环境）
2. **`.env` 文件**（自动生成，用于开发环境）
3. **自动生成**（首次启动时）

#### 首次启动

首次启动时，如果未设置环境变量，系统会：
- 自动生成一个安全的随机密钥（32字节，Base64编码）
- 保存到项目根目录的 `.env` 文件
- 设置文件权限为 `0600`（仅所有者可读写）
- 在日志中提示备份 `.env` 文件

```
[Encryption] No master key found, generating new key...
[Encryption] New master key generated and saved to .env
[Encryption] ⚠️  IMPORTANT: Backup this .env file! Losing it means permanent data loss.
```

#### ⚠️ 重要：备份 .env 文件

**`.env` 文件包含加密密钥，丢失后将无法解密数据库中的敏感信息！**

```bash
# 备份 .env 文件到安全位置
cp .env .env.backup

# 或者使用环境变量（推荐生产环境）
export BITWARDEN_BACKUP_MASTER_KEY="your-key-from-env-file"
```

#### 生产环境部署

**推荐使用环境变量而非 .env 文件：**

```bash
# 方式 1: 直接设置环境变量
export BITWARDEN_BACKUP_MASTER_KEY="your-generated-key-here"

# 方式 2: 从 .env 文件读取并设置
export BITWARDEN_BACKUP_MASTER_KEY=$(grep BITWARDEN_BACKUP_MASTER_KEY .env | cut -d '=' -f2 | tr -d '"')
```

**Docker Compose 部署：**

```yaml
# docker-compose.yml
services:
  bitwarden-backup:
    environment:
      - BITWARDEN_BACKUP_MASTER_KEY=${BITWARDEN_BACKUP_MASTER_KEY}
```

#### 手动生成密钥（可选）

如果需要手动生成密钥：

```bash
# 生成密钥
openssl rand -base64 32

# 设置环境变量
export BITWARDEN_BACKUP_MASTER_KEY="your-generated-key-here"
```

## 使用指南

### 1. 配置源服务器

在"源服务器"页面添加 Bitwarden 服务器配置：

- **名称**: 服务器标识名称
- **服务器 URL**: Bitwarden 服务器地址（如 `https://vault.bitwarden.com`）
- **Client ID / Client Secret**: API 密钥（在 Bitwarden 设置中获取）
- **Master Password**: 主密码（用于解锁密码库导出数据）

> **注意**: 即使使用 API Key 登录，Bitwarden CLI 仍需要主密码来解锁密码库并导出数据。

> **🔒 安全**: 所有敏感凭证（Client ID、Client Secret、Master Password）都使用 AES-256-GCM 加密存储。详见 [SECURITY.md](SECURITY.md)

### 2. 配置备份目标

支持三种备份目标类型：

- **本地存储**: 备份到本地目录
- **WebDAV**: 备份到 WebDAV 服务器（如 Nextcloud）
- **目标服务器**: 导入到另一个 Bitwarden 服务器

### 3. 创建备份任务

配置定时备份任务，支持 Cron 表达式：

```
秒 分 时 日 月 周
0  0  2  *  *  *   # 每天凌晨 2 点
0  0  */6 *  *  *  # 每 6 小时
0  30 1  *  *  1   # 每周一凌晨 1:30
```

## 目录结构

```
bitwarden-backup/
├── data/           # 数据库文件
├── backups/        # 本地备份文件
└── .tmp/           # 临时文件（自动清理）
```

## License

MIT
