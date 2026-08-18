package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadyHandler(t *testing.T) {
	t.Run("初期状態では 200 を返す", func(t *testing.T) {
		r := newReadiness()
		rec := httptest.NewRecorder()

		r.handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("停止開始後は 503 を返す", func(t *testing.T) {
		r := newReadiness()
		r.markShuttingDown()
		rec := httptest.NewRecorder()

		r.handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ready", nil))

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
	})
}
