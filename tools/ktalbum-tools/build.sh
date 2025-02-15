#!/bin/bash

# 版本号
VERSION="v0.1.0"

# 构建前端
echo "构建前端..."
cd web/frontend
npm install

# 直接构建，暂时跳过类型检查
npm run build
if [ $? -ne 0 ]; then
    echo "前端构建失败"
    exit 1
fi

# 确保构建目录存在
if [ ! -d "dist" ]; then
    echo "前端构建目录不存在"
    exit 1
fi

cd ../..

# 创建发布目录
mkdir -p release

# 构建各平台版本
echo "构建各平台版本..."

# Windows (64-bit, x86_64)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o release/ktalbum-tools-${VERSION}-windows-amd64.exe

# Windows (64-bit, ARM64)
GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o release/ktalbum-tools-${VERSION}-windows-arm64.exe

# macOS (64-bit, Intel)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o release/ktalbum-tools-${VERSION}-darwin-amd64

# macOS (64-bit, Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o release/ktalbum-tools-${VERSION}-darwin-arm64

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/ktalbum-tools-${VERSION}-linux-amd64

# 为 Unix-like 系统添加执行权限
chmod +x release/ktalbum-tools-${VERSION}-darwin-*
chmod +x release/ktalbum-tools-${VERSION}-linux-*

echo "构建完成！" 