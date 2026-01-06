#!/bin/bash

# Bitwarden Backup 启动脚本

set -e

echo "==================================="
echo "  Bitwarden Backup Manager"
echo "==================================="
echo ""

# 创建必要的目录
mkdir -p data backups

# 检查是否存在编译好的二进制文件
if [ -f "./bitwarden-backup" ]; then
    echo "启动服务..."
    ./bitwarden-backup
elif [ -f "./bitwarden-backup.exe" ]; then
    echo "启动服务..."
    ./bitwarden-backup.exe
else
    # 尝试编译
    if command -v go &> /dev/null; then
        echo "编译项目..."
        go build -o bitwarden-backup ./cmd/server
        echo "启动服务..."
        ./bitwarden-backup
    else
        echo "错误: 未找到可执行文件，且 Go 未安装"
        echo ""
        echo "请选择以下方式之一:"
        echo "  1. 安装 Go 后运行: go build -o bitwarden-backup ./cmd/server && ./bitwarden-backup"
        echo "  2. 使用 Docker: docker compose up -d"
        exit 1
    fi
fi
