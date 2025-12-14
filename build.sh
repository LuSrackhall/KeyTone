#!/bin/bash
# =============================================================================
# 统一构建入口脚本
# 用法: ./build.sh <platform>
# 支持平台: win | mac | linux
# =============================================================================

set -euo pipefail  # 黄金三件套，严格模式
IFS=$'\n\t'        # 安全的字段分隔符

cd sdk

# 颜色输出（好看 + 专业）
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info()    { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC} $*"; }
error()   { echo -e "${RED}[ERROR]${NC} $*" >&2; }
die()     { error "$*"; exit 1; }

# ============ 参数解析 ============
show_help() {
    cat << EOF
用法: ./build.sh <平台>

支持的平台:
  win       构建 Windows 版本（默认）
  mac       构建 macOS 版本
  linux     构建 Linux 版本

示例:
  ./build.sh win
  ./build.sh mac
  ./build.sh linux
  ./build.sh -h        显示此帮助

提示: 你也可以直接运行 ./build.sh（不带参数）默认构建 Windows
EOF
}

PLATFORM="win"  # 默认平台

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        win|mac|linux)
            PLATFORM="$1"
            shift
            ;;
        *)
            error "未知参数: $1"
            show_help
            exit 1
            ;;
    esac
done

info "目标平台: $PLATFORM"

# ============ 加载密钥环境（关键！）============
info "正在加载密钥环境..."
if source ./setup_build_env.sh; then
    info "密钥环境加载成功"
else
    die "密钥环境加载失败，请检查 private_keys.env 是否存在并已填写"
fi

# ============ 执行构建 ============

case "$PLATFORM" in
    win)
        info "开始构建 Windows 版本..."
        make win
        ;;
    mac)
        info "开始构建 macOS 版本..."
        make mac
        ;;
    linux)
        info "开始构建 Linux 版本..."
        make linux
        ;;
    *)
        die "不支持的平台: $PLATFORM"
        ;;
esac

# ============ 成功提示 ============
info "恭喜！$PLATFORM 构建完成！"
echo
echo "输出目录示例（根据你的 Makefile 实际为准）:"
case "$PLATFORM" in
    win)   echo "   ./dist/KeyToneSetup.exe" ;;
    mac)   echo "   ./dist/KeyTone.app" ;;
    linux) echo "   ./dist/keytone-linux-x64" ;;
esac
echo