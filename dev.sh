#!/bin/bash

# Vue 前端构建和 Go 后端启动脚本

set -e

echo "🚀 开始构建 Vue 前端..."

# 进入前端目录
cd web/frontend

# 检查 node_modules 是否存在
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    npm install
fi

# 构建前端
echo "🔨 构建前端..."
npm run build

cd ../..

echo "✅ 前端构建完成！"
echo ""
echo "🔧 启动 Go 后端服务器..."
echo ""

# 启动 Go 服务器
go run ./cmd/server
