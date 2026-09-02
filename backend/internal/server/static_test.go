package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// テスト用の擬似的なビルド成果物。
// 本番では embed.FS を渡すが、単体テストでは実ファイルを不要にするため
// fstest.MapFS を注入する。
func testAssets() fstest.MapFS {
	return fstest.MapFS{
		"index.html":    {Data: []byte("<!doctype html>INDEX")},
		"assets/app.js": {Data: []byte("console.log(1)")},
		"favicon.svg":   {Data: []byte("<svg/>")},
	}
}

func TestSPAHandler(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantBody string
	}{
		{"ルートは index.html を返す", "/", "<!doctype html>INDEX"},
		{"実在する静的ファイルはそのまま返す", "/assets/app.js", "console.log(1)"},
		{"実在するアセットはそのまま返す", "/favicon.svg", "<svg/>"},
		{"実在しないパスは index.html に fallback する", "/about", "<!doctype html>INDEX"},
		{"深いパスも fallback する", "/works/detail/1", "<!doctype html>INDEX"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := spaHandler(testAssets())
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
			body, _ := io.ReadAll(rec.Body)
			if string(body) != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
