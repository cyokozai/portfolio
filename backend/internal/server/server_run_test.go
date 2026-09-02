package server

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 停止時は「readiness を 503 にする」→「HTTP サーバを閉じる」の順でなければ
// ならない。逆順だと Service がまだ Pod を宛先に含んだまま接続が切れ 502 になる。
func TestServer_ReadinessGoesDownBeforeShutdown(t *testing.T) {
	s := New(testAssets())
	s.drainDelay = 0

	var codeAtShutdown int
	s.shutdownFunc = func(context.Context) error {
		rec := httptest.NewRecorder()
		s.ready.handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))
		codeAtShutdown = rec.Code
		return nil
	}

	if err := s.shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}

	if codeAtShutdown != http.StatusServiceUnavailable {
		t.Errorf("Shutdown 実行時点の /ready = %d, want %d（先に readiness を落とすこと）",
			codeAtShutdown, http.StatusServiceUnavailable)
	}
}

func TestServer_ServeStopsOnContextCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	s := New(testAssets())
	s.drainDelay = 0

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() { errCh <- s.Serve(ctx, ln) }()

	base := "http://" + ln.Addr().String()
	resp, err := http.Get(base + "/ready")
	if err != nil {
		t.Fatalf("起動後の /ready に到達できない: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("/ready = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Serve = %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("ctx キャンセル後も Serve が終了しない")
	}
}
