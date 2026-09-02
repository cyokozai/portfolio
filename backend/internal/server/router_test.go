package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		wantStatus   int
		wantContains string
	}{
		{"health は 200", "/health", http.StatusOK, `"status":"ok"`},
		{"ready は 200", "/ready", http.StatusOK, `"status":"ready"`},
		{"ルートは SPA を返す", "/", http.StatusOK, "INDEX"},
		{"SPA のルートは fallback する", "/about", http.StatusOK, "INDEX"},
		{"静的アセットは配信される", "/assets/app.js", http.StatusOK, "console.log(1)"},
		// ここが本質: /api 配下は SPA に fallback せず 404 を返すこと
		{"未定義の API は 404", "/api/unknown", http.StatusNotFound, ""},
		{"API のルートも 404", "/api/", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewRouter(testAssets(), newReadiness())
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			body, _ := io.ReadAll(rec.Body)
			if tt.wantContains != "" && !strings.Contains(string(body), tt.wantContains) {
				t.Errorf("body = %q, want contains %q", body, tt.wantContains)
			}
			if tt.wantStatus == http.StatusNotFound && strings.Contains(string(body), "INDEX") {
				t.Errorf("API パスが SPA に fallback している: %q", body)
			}
		})
	}
}
