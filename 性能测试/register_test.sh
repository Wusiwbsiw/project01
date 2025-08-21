#!/bin/bash

# ==============================================================================
# HTTP API 注册接口并发压力测试脚本 (最终修正版)
#
# 修正: 使用 `timeout` 命令来确保目标生成器会在测试结束后自动停止，
#       从而避免管道阻塞导致脚本挂起的问题。
# ==============================================================================

# --- 可配置的参数 ---
TARGET_URL="https://localhost:8443/api/user/register"
RATE=${1:-50}
DURATION=${2:-"10s"}

# --- 核心逻辑 ---

echo "🚀 开始对注册接口进行压力测试..."
echo "-------------------------------------"
echo "  目标 URL : $TARGET_URL"
echo "  请求频率 (Rate) : $RATE req/s"
echo "  持续时间 (Duration) : $DURATION"
echo "-------------------------------------"
echo ""

# 使用 timeout 命令来确保我们的“目标生成器”会在测试结束后自动停止
# bash -c '...' 用于将整个 while 循环作为单个命令传递给 timeout
timeout "$DURATION" \
bash -c 'i=0; while true; do
  i=$((i+1));
  jq -cn \
    --arg method "POST" \
    --arg url "'"$TARGET_URL"'" \
    --arg user "testuser_$i" \
    --arg pass "password123" \
    '\''{"method": $method, "url": $url, "header": {"Content-Type": ["application/json"]}, "body": ({"username": $user, "password": $pass} | @base64)}'\'';
done' \
| vegeta attack -format=json -rate="$RATE" -duration="$DURATION" -insecure \
| vegeta report

echo ""
echo "✅ 测试完成。"
