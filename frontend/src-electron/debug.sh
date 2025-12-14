#!/bin/bash

# 切换到 SDK 目录(使用 || exit 确保目录切换失败时脚本会终止)
cd ../sdk || exit

# 执行 Go 构建命令
# -ldflags "$EXTRA_LDFLAGS": 注入环境变量中定义的额外链接参数（如混淆后的密钥）
# 如果 EXTRA_LDFLAGS 为空，则不注入任何额外参数，使用代码中的默认值。
go build -ldflags "$EXTRA_LDFLAGS" -o ../frontend/src-electron/sdk-debug/KeyTone.exe
