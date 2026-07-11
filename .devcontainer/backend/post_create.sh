#!/usr/bin/env bash
set -e

echo "==> [postCreate] portfolio-backend セットアップ開始"

cd /workspace/backend

if [ -f go.mod ]; then
  echo "==> go.mod を検出。依存を取得します (go mod download)"
  go mod download
  go mod verify
else
  echo "==> go.mod が見つかりません。次を実行してモジュールを初期化してください:"
  echo ""
  echo "      cd /workspace/backend"
  echo "      go mod init github.com/cyokozai/portfolio/backend"
  echo ""
fi

echo "==> Go: $(go version)"
echo "==> tools: gopls / dlv / staticcheck を同梱済み"
echo "==> [postCreate] 完了"
