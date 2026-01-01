#!/bin/bash

# setup_build_env.sh (ktalbum-tools)
# ======================================================================================
# 目标：为 ktalbum-tools 构建生成 EXTRA_LDFLAGS（-ldflags -X ...），用于注入本地私有密钥。
# - 读取 private_keys.env（优先使用 tools/ktalbum-tools/private_keys.env，其次复用 sdk/private_keys.env）
# - 调用 tools/key-obfuscator 生成 XOR 混淆后的 hex
# - 拼接为 Go 链接器参数
#
# 用法：
#   source ./setup_build_env.sh
# 或：
#   eval $(./setup_build_env.sh)
# ======================================================================================

# =================配置区域=================

KEYS_FILE_LOCAL="private_keys.env"
KEYS_FILE_SDK="../../sdk/private_keys.env"

# 选择优先的 keys 文件
if [ -f "$KEYS_FILE_LOCAL" ]; then
  KEYS_FILE="$KEYS_FILE_LOCAL"
else
  KEYS_FILE="$KEYS_FILE_SDK"
fi

OBFUSCATOR_TOOL="../key-obfuscator/main.go"

# 严格模式开关（可选）
STRICT_KEYS=${KEYTONE_STRICT_KEYS:-0}

# 格式： "环境变量中的键名:Go代码中的变量全路径"
KEYS_TO_PROCESS=(
  "KEY_V1:ktalbum-tools/utils.KeytoneEncryptKeyV1"
  "KEY_V2:ktalbum-tools/utils.KeytoneEncryptKeyV2"
)

# =================逻辑区域=================

if [ ! -f "$KEYS_FILE" ]; then
  if [ "$STRICT_KEYS" = "1" ]; then
    echo "错误: 未找到私钥文件 $KEYS_FILE" >&2
    echo "请复制 sdk/private_keys.template.env 为 sdk/private_keys.env 或 tools/ktalbum-tools/private_keys.env 并填入密钥。" >&2
    return 1 2>/dev/null || exit 1
  fi

  echo "警告: 未找到私钥文件 $KEYS_FILE，将跳过密钥注入并使用源码默认密钥。" >&2
  export EXTRA_LDFLAGS=""
  if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "export EXTRA_LDFLAGS=\"\""
  fi
  return 0 2>/dev/null || exit 0
fi

if [ ! -f "$OBFUSCATOR_TOOL" ]; then
  if [ "$STRICT_KEYS" = "1" ]; then
    echo "错误: 未找到混淆工具源码 $OBFUSCATOR_TOOL" >&2
    return 1 2>/dev/null || exit 1
  fi

  echo "警告: 未找到混淆工具源码 $OBFUSCATOR_TOOL，将跳过密钥注入并使用源码默认密钥。" >&2
  export EXTRA_LDFLAGS=""
  if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "export EXTRA_LDFLAGS=\"\""
  fi
  return 0 2>/dev/null || exit 0
fi

LDFLAGS=""

# 在上层启用 `set -u` 时，确保这些变量始终已定义
KEY_NAME=""
GO_VAR=""
PLAINTEXT_KEY=""
OBFUSCATED_VAL=""

echo "正在处理 ktalbum-tools 密钥混淆..." >&2

for entry in "${KEYS_TO_PROCESS[@]}"; do
  KEY_NAME=$(echo "$entry" | cut -d':' -f1)
  GO_VAR=$(echo "$entry" | cut -d':' -f2)

  # 兼容上层可能启用 `set -euo pipefail` 的场景：
  # grep 找不到时必须返回 0，避免脚本被直接中断。
  PLAINTEXT_KEY=$( (grep "^${KEY_NAME}=" "${KEYS_FILE}" || true) | head -n 1 | cut -d'"' -f2 )

  if [ -z "$PLAINTEXT_KEY" ]; then
    if [ "$STRICT_KEYS" = "1" ]; then
      echo "错误: 在 ${KEYS_FILE} 中未找到 ${KEY_NAME}" >&2
      return 1 2>/dev/null || exit 1
    fi

    echo "警告: 在 ${KEYS_FILE} 中未找到 ${KEY_NAME}，将跳过该项注入并使用源码默认值。" >&2
    continue
  fi

  if ! OBFUSCATED_VAL=$(go run "${OBFUSCATOR_TOOL}" -key "${PLAINTEXT_KEY}"); then
    if [ "$STRICT_KEYS" = "1" ]; then
      echo "错误: 密钥 ${KEY_NAME} 混淆失败" >&2
      return 1 2>/dev/null || exit 1
    fi

    echo "警告: 密钥 ${KEY_NAME} 混淆失败，将跳过该项注入并使用源码默认值。" >&2
    continue
  fi

  LDFLAGS="$LDFLAGS -X '${GO_VAR}=${OBFUSCATED_VAL}'"
 done

export EXTRA_LDFLAGS="$LDFLAGS"

if [ -z "$LDFLAGS" ]; then
  echo "成功！未注入任何密钥，EXTRA_LDFLAGS 为空，将使用源码默认密钥（ktalbum-tools）。" >&2
else
  echo "成功！已设置 EXTRA_LDFLAGS 环境变量（ktalbum-tools）。" >&2
fi

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  echo "export EXTRA_LDFLAGS=\"$LDFLAGS\""
fi
