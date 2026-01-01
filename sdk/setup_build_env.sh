#!/bin/bash

# setup_build_env.sh
# ======================================================================================
# 脚本功能说明：
# 此脚本用于在本地开发或构建前，自动读取私钥文件，调用混淆工具生成混淆后的密钥，
# 并将其组装成 Go 编译器需要的链接参数 (EXTRA_LDFLAGS)。
#
# 为什么需要这个脚本？
# 1. 安全性：私钥以明文存储在本地文件（不提交到 git），构建时注入的是混淆后的 Hex 值。
# 2. 自动化：免去手动运行混淆工具和复制粘贴 Hex 值的繁琐步骤。
# 3. 兼容性：生成的 EXTRA_LDFLAGS 变量可以被 Makefile 和 debug.sh 脚本直接使用。
# ======================================================================================

# 使用方法：
# 方式一（推荐）：在当前终端加载环境变量
#   source ./setup_build_env.sh
#   > 这种方式会将 EXTRA_LDFLAGS 导出到当前终端会话中。
#   > 当您关闭终端后，该环境变量会自动失效，不会污染您的全局系统环境。
#
# 方式二（高级）：仅获取 export 命令
#   ./setup_build_env.sh
#   > 这种方式只会打印出 export 命令，不会改变当前环境。
#   > 您可以配合 eval 使用：eval $(./setup_build_env.sh)

# =================配置区域=================

# 1. 私钥文件路径
# 为什么使用 .env 格式？
# .env 是业界通用的"键=值"配置文件格式，易于阅读，也方便脚本解析。
# 虽然它叫 .env，但我们并不直接将其加载为系统环境变量，而是作为普通文本文件读取。
KEYS_FILE="private_keys.env"

# 2. 混淆工具源码路径
OBFUSCATOR_TOOL="../tools/key-obfuscator/main.go"

# 4. 严格模式开关（可选）
# 默认（0）：任意 KEY_* 缺失都不会导致失败，自动跳过注入并使用代码默认值
# 严格（1）：一旦缺少私钥文件或任意 KEY_*，立即失败
STRICT_KEYS=${KEYTONE_STRICT_KEYS:-0}

# 3. 定义需要处理的密钥列表
# 格式： "环境变量中的键名:Go代码中的变量全路径"
# 如果后续需要新增密钥，只需在此数组中追加一行即可。
# 格式说明：
#   KEY_NAME  : 在 private_keys.env 文件中的键名 (如 KEY_F)
#   GO_VAR    : 在 Go 代码中接收注入的变量全路径 (如 KeyTone/signature.KeyToneAuthRequestEncryptionKeyF)
KEYS_TO_PROCESS=(
    "KEY_F:KeyTone/signature.KeyToneAuthRequestEncryptionKeyF"
    "KEY_K:KeyTone/signature.KeyToneAuthRequestEncryptionKeyK"
    "KEY_Y:KeyTone/signature.KeyToneAuthGrantEncryptionKeyY"
    "KEY_N:KeyTone/signature.KeyToneAuthGrantEncryptionKeyN"

    # ===== 对称密钥（历史遗留，现统一适配构建注入） =====
    "KEY_A:KeyTone/signature.KeyToneSignatureEncryptionKeyA"
    "KEY_B:KeyTone/signature.KeyToneSignatureEncryptionKeyB"
    "KEY_V1:KeyTone/server.KeytoneEncryptKeyV1"
    "KEY_V2:KeyTone/server.KeytoneEncryptKeyV2"
    "FIXED_SECRET:KeyTone/audioPackage/enc.FixedSecret"
    "KEY_ALBUM_SIG:KeyTone/signature.KeyToneAlbumSignatureEncryptionKey"
    # 示例：新增密钥时，取消注释并修改下行
    # "KEY_NEW:KeyTone/signature.KeyToneNewKey"
)

# =================逻辑区域=================

# 检查私钥文件是否存在
if [ ! -f "$KEYS_FILE" ]; then
    if [ "$STRICT_KEYS" = "1" ]; then
        echo "错误: 未找到私钥文件 $KEYS_FILE" >&2
        echo "请复制 private_keys.template.env 为 $KEYS_FILE 并填入您的密钥。" >&2
        return 1 2>/dev/null || exit 1
    fi

    echo "警告: 未找到私钥文件 $KEYS_FILE，将跳过密钥注入并使用源码默认密钥。" >&2
    export EXTRA_LDFLAGS=""
    if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
        echo "export EXTRA_LDFLAGS=\"\""
    fi
    return 0 2>/dev/null || exit 0
fi

# 检查混淆工具是否存在
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

# 初始化 LDFLAGS 字符串
LDFLAGS=""

# 在上层启用 `set -u` 时，确保这些变量始终已定义
KEY_NAME=""
GO_VAR=""
PLAINTEXT_KEY=""
OBFUSCATED_VAL=""

echo "正在处理密钥混淆..." >&2

# 遍历密钥列表进行处理
for entry in "${KEYS_TO_PROCESS[@]}"; do
    # 使用 cut 命令分割字符串
    # -d':' 指定冒号为分隔符
    # -f1 取第一部分(键名), -f2 取第二部分(Go变量路径)
    KEY_NAME=$(echo "$entry" | cut -d':' -f1)
    GO_VAR=$(echo "$entry" | cut -d':' -f2)

    # 从文件中读取明文密钥
    # grep "^$KEY_NAME=" : 查找以 KEY_NAME= 开头的行
    # cut -d'"' -f2      : 以双引号为分隔符，提取中间的内容(即密钥值)
    # 注意：上层脚本（dev.sh/build.sh）启用了 `set -euo pipefail`
    # 这里必须避免 grep 找不到时导致整条 pipeline 非 0 直接中断。
    PLAINTEXT_KEY=$( (grep "^${KEY_NAME}=" "${KEYS_FILE}" || true) | head -n 1 | cut -d'"' -f2 )

    # 检查是否成功读取到密钥
    if [ -z "$PLAINTEXT_KEY" ]; then
        if [ "$STRICT_KEYS" = "1" ]; then
            echo "错误: 在 ${KEYS_FILE} 中未找到 ${KEY_NAME}" >&2
            return 1 2>/dev/null || exit 1
        fi

        echo "警告: 在 ${KEYS_FILE} 中未找到 ${KEY_NAME}，将跳过该项注入并使用源码默认值。" >&2
        continue
    fi

    # 调用 Go 工具生成混淆后的 Hex 字符串
    # go run ... : 直接运行 Go 源码，无需预先编译
    # -key ...   : 传递明文密钥作为参数
    # $()        : 命令替换，将命令的输出结果赋值给变量
    # 同理，go run 失败也不能在默认模式下中断整个 dev/build
    if ! OBFUSCATED_VAL=$(go run "${OBFUSCATOR_TOOL}" -key "${PLAINTEXT_KEY}"); then
        if [ "$STRICT_KEYS" = "1" ]; then
            echo "错误: 密钥 ${KEY_NAME} 混淆失败" >&2
            return 1 2>/dev/null || exit 1
        fi

        echo "警告: 密钥 ${KEY_NAME} 混淆失败，将跳过该项注入并使用源码默认值。" >&2
        continue
    fi

    # 拼接到 LDFLAGS 字符串中
    # -X ... : Go 链接器参数，用于设置变量值
    LDFLAGS="$LDFLAGS -X '${GO_VAR}=${OBFUSCATED_VAL}'"
done

# 导出环境变量
# export 命令将变量设置为环境变量，使其对当前 shell 及其子进程可见
export EXTRA_LDFLAGS="$LDFLAGS"

if [ -z "$LDFLAGS" ]; then
    echo "成功！未注入任何密钥，EXTRA_LDFLAGS 为空，将使用源码默认密钥。" >&2
else
    echo "成功！已设置 EXTRA_LDFLAGS 环境变量（包含混淆后的密钥）。" >&2
fi
echo "您现在可以运行 'make <target>' 或启动 electron 开发环境。" >&2
echo "注意：此环境变量仅在当前终端会话有效，关闭终端后自动失效，不会污染全局环境。" >&2

# 兼容 eval 用法
# 如果脚本是被直接执行的(而不是 source)，打印 export 命令供 eval 使用
# ${BASH_SOURCE[0]} 表示当前脚本的路径
# $0 表示当前执行的命令
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "export EXTRA_LDFLAGS=\"$LDFLAGS\""
fi
