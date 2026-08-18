package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/cyokozai/portfolio/backend/internal/assets"
	"github.com/cyokozai/portfolio/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("サーバが異常終了しました", "error", err)
		os.Exit(1)
	}
	slog.Info("サーバを正常に停止しました")
}

func run() error {
	// SIGTERM は k8s が Pod を停止するときに送る。
	// ここで受けて Server 側の graceful shutdown に繋ぐ。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	slog.Info("待ち受けを開始します", "addr", addr)
	return server.New(assets.FS()).Run(ctx, addr)
}
