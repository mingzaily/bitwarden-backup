#!/bin/bash

# Bitwarden Backup 启动脚本

echo "==================================="
echo "  Bitwarden Backup Manager"
echo "==================================="
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: Docker 未安装"
    echo "请先安装 Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "错误: Docker Compose 未安装"
    echo "请先安装 Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

# 创建必要的目录
echo "创建数据目录..."
mkdir -p data backups

# 构建并启动容器
echo "构建 Docker 镜像..."
docker-compose build

echo "启动服务..."
docker-compose up -d

echo ""
echo "✅ 服务启动成功！"
echo ""
echo "访问地址: http://localhost:8080"
echo ""
echo "常用命令:"
echo "  查看日志: docker-compose logs -f"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart"
echo ""
