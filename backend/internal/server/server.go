package server

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"time"
)

// drainDelay の既定値。readiness を 503 にしてから接続を閉じるまでの待ち時間。
// kube-proxy / Ingress がエンドポイントの除外を反映するのを待つためのもの。
const defaultDrainDelay = 3 * time.Second

// Server は SPA と API を単一プロセスで配信する HTTP サーバ。
type Server struct {
	ready      *readiness
	handler    http.Handler
	drainDelay time.Duration
	httpServer *http.Server

	// shutdownFunc は停止処理の実体。Serve が実サーバの Shutdown を設定する。
	// テストでは差し替えて停止順序を検証する。
	shutdownFunc func(context.Context) error
}

func New(assets fs.FS) *Server {
	ready := newReadiness()
	return &Server{
		ready:      ready,
		handler:    NewRouter(assets, ready),
		drainDelay: defaultDrainDelay,
	}
}

// Run は addr で listen して Serve する。
func (s *Server) Run(ctx context.Context, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(ctx, ln)
}

// Serve は ln で待ち受け、ctx がキャンセルされたら graceful に停止する。
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	s.httpServer = &http.Server{
		Handler:           s.handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	s.shutdownFunc = s.httpServer.Shutdown

	serveErr := make(chan error, 1)
	go func() {
		err := s.httpServer.Serve(ln)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveErr <- err
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
	defer cancel()
	if err := s.shutdown(shutdownCtx); err != nil {
		return err
	}
	return <-serveErr
}

// shutdown は readiness を落としてから HTTP サーバを閉じる。
// この順序が逆になると、Service がまだこの Pod を宛先に含んでいる間に
// 接続が切れて 502 になる。
func (s *Server) shutdown(ctx context.Context) error {
	s.ready.markShuttingDown()

	if s.drainDelay > 0 {
		select {
		case <-time.After(s.drainDelay):
		case <-ctx.Done():
		}
	}

	if s.shutdownFunc == nil {
		return nil
	}
	return s.shutdownFunc(ctx)
}
