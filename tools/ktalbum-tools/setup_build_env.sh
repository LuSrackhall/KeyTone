#!/bin/bash

# tools/ktalbum-tools/setup_build_env.sh
# ======================================================================================
# 脚本功能说明：
# 该脚本用于 ktalbum-tools 的本地构建（个人调试/查看用），支持与 SDK 相同的“构建时密钥注入”机制。
#
# 设计目标：
# - 开源构建：不需要任何私钥文件即可工作（使用默认开源密钥）
# - 私有构建：可读取 sdk/private_keys.env（或本地 private_keys.env），生成混淆后的 Hex，并注入到 ktalbum-tools
# - 工具兼容性：ktalbum-tools 运行时会同时尝试“注入密钥”与“开源默认密钥”进行解密，以兼容两类产物
# ======================================================================================

# 使用方法：
# 方式一（推荐）：在当前终端加载环境变量
#   source ./setup_build_env.sh
#
# 方式二：仅获取 export 命令（可配合 eval）
#   ./setup_build_env.sh
#   eval $(./setup_build_env.sh)

# =================配置区域=================

# 1. 私钥文件路径
# 优先复用 SDK 的私钥文件（不提交到 git）
# 可通过环境变量 KEYS_FILE 覆盖
DEFAULT_KEYS_FILE="../../sdk/private_keys.env"
KEYS_FILE="${KEYS_FILE:-$DEFAULT_KEYS_FILE}"

# 如果默认路径不存在，尝试使用当前目录的 private_keys.env
if [ ! -f "$KEYS_FILE" ]; then
  if [ -f "private_keys.env" ]; then
    KEYS_FILE="private_keys.env"
  fi
fi

# 2. 混淆工具源码路径（与 SDK 共用同一份）
OBFUSCATOR_TOOL="../key-obfuscator/main.go"

# 3. 定义需要处理的密钥列表
# ktalbum-tools 只需要专辑导出文件 XOR key（v1/v2）即可解密 .ktalbum
# 格式： "环境变量中的键名:Go代码中的变量全路径"
KEYS_TO_PROCESS=(
  "KEY_ALBUM_EXPORT_V1:ktalbum-tools/utils.KeytoneEncryptKeyV1"
  "KEY_ALBUM_EXPORT_V2:ktalbum-tools/utils.KeytoneEncryptKeyV2"
)

# =================逻辑区域=================

# 检查混淆工具是否存在
if [ ! -f "$OBFUSCATOR_TOOL" ]; then
  echo "错误: 未找到混淆工具源码 $OBFUSCATOR_TOOL" >&2
  return 1 2>/dev/null || exit 1
fi

# 如果没有私钥文件，允许继续（开源构建不需要 EXTRA_LDFLAGS）
if [ ! -f "$KEYS_FILE" ]; then
  echo "提示: 未找到私钥文件 $KEYS_FILE，将不设置 EXTRA_LDFLAGS（开源默认密钥模式）。" >&2
  export EXTRA_LDFLAGS=""
  if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "export EXTRA_LDFLAGS=\"\""
  fi
  exit 0
fi

# 初始化 LDFLAGS 字符串
LDFLAGS=""

echo "正在处理 ktalbum-tools 密钥混淆..." >&2

for entry in "${KEYS_TO_PROCESS[@]}"; do
  KEY_NAME=$(echo "$entry" | cut -d':' -f1)
  GO_VAR=$(echo "$entry" | cut -d':' -f2)

  PLAINTEXT_KEY=$(grep "^$KEY_NAME=" "$KEYS_FILE" | cut -d'"' -f2)
  if [ -z "$PLAINTEXT_KEY" ]; then
    echo "错误: 在 $KEYS_FILE 中未找到 $KEY_NAME" >&2
    return 1 2>/dev/null || exit 1
  fi

  OBFUSCATED_VAL=$(go run "$OBFUSCATOR_TOOL" -key "$PLAINTEXT_KEY")
  if [ $? -ne 0 ]; then
    echo "错误: 密钥 $KEY_NAME 混淆失败" >&2
    return 1 2>/dev/null || exit 1
  fi

  LDFLAGS="$LDFLAGS -X '$GO_VAR=$OBFUSCATED_VAL'"
done

export EXTRA_LDFLAGS="$LDFLAGS"

echo "成功！已设置 EXTRA_LDFLAGS（ktalbum-tools）。" >&2

echo "示例构建：" >&2

echo "  go build -ldflags \"$EXTRA_LDFLAGS\" ./..." >&2

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "export EXTRA_LDFLAGS=\"$LDFLAGS\""
fi
