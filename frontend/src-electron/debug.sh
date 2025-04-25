#!/bin/bash

# 切换到 SDK 目录(使用 || exit 确保目录切换失败时脚本会终止)
cd ../sdk || exit

# 执行 Go 构建命令
go build -o ../frontend/src-electron/sdk-debug/KeyTone.exe
