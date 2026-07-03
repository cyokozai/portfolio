#!/usr/bin/env bash
set -e

echo "==> [postCreate] portfolio-frontend セットアップ開始"

if [ -f package.json ]; then
  echo "==> package.json を検出。依存関係をインストールします (npm install)"
  npm install
else
  echo "==> package.json が見つかりません（まだ雛形が無い状態）。"
  echo "    コンテナ内のターミナルで、次を実行して Svelte 雛形を作成してください:"
  echo ""
  echo "      npm create vite@latest . -- --template svelte"
  echo "      npm install"
  echo ""
  echo "    （'.' は現在のフォルダに作る指定。既存ファイルがあると警告が出るので"
  echo "      'Ignore files and continue' を選んでください）"
fi

echo "==> [postCreate] 完了"
