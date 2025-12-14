#!/bin/bash
set -euo pipefail

cd sdk

# 加载密钥环境（失败直接退出，不需要额外提示）
source ./setup_build_env.sh || exit 1

cd ../frontend
quasar dev -m electron