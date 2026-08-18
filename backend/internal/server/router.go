package server

import (
	"encoding/json"
	"io/fs"
	"net/http"
)

// NewRouter はアプリケーション全体のルーティングを組み立てる。
//
//	/health  liveness probe
//	/ready   readiness probe（停止開始後は 503）
//	/api/*   コンテンツ API。未定義パスは 404（SPA へ fallback させない）
//	/*       SPA。実ファイルが無ければ index.html を返す
func NewRouter(assets fs.FS, ready *readiness) http.Handler {
	mux := http.NewServeMux()

	// probe 系。/health は liveness なので停止処理中も 200 を返し続ける。
	mux.Handle("GET /health", healthHandler())
	mux.Handle("GET /ready", ready.handler())

	// コンテンツ API
	mux.Handle("GET /api/profile", profileHandler())

	// ServeMux はより具体的なパターンを優先する。この catch-all があることで
	// 未定義の /api/* が SPA の fallback に流れず 404 になる。
	mux.Handle("/api/", apiNotFoundHandler())

	mux.Handle("/", spaHandler(assets))

	return mux
}

// apiNotFoundHandler は未定義の API パスに JSON で 404 を返す。
func apiNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
	})
}
