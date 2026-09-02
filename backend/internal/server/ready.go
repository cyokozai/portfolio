package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
)

// readiness は readiness probe の応答状態を保持する。
// SIGTERM 受信時にまず 503 へ反転させ、Service のエンドポイントから
// 外れてから接続を閉じるために使う。
type readiness struct {
	shuttingDown atomic.Bool
}

func newReadiness() *readiness {
	return &readiness{}
}

// markShuttingDown は以降の readiness probe を 503 にする。
func (r *readiness) markShuttingDown() {
	r.shuttingDown.Store(true)
}

func (r *readiness) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if r.shuttingDown.Load() {
			http.Error(w, "shutting down", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})
}
