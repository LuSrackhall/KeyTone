#!/bin/bash

# 版本号
VERSION="v0.1.0"

# 创建发布目录
mkdir -p release

# Windows (64-bit, x86_64)
GOOS=windows GOARCH=amd64 go build -o release/ktalbum-tools-${VERSION}-windows-amd64.exe

# Windows (64-bit, ARM64)
GOOS=windows GOARCH=arm64 go build -o release/ktalbum-tools-${VERSION}-windows-arm64.exe

# macOS (64-bit, Intel)
GOOS=darwin GOARCH=amd64 go build -o release/ktalbum-tools-${VERSION}-darwin-amd64

# macOS (64-bit, Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o release/ktalbum-tools-${VERSION}-darwin-arm64

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o release/ktalbum-tools-${VERSION}-linux-amd64

# 为 Unix-like 系统添加执行权限
chmod +x release/ktalbum-tools-${VERSION}-darwin-*
chmod +x release/ktalbum-tools-${VERSION}-linux-*

echo "构建完成！" 