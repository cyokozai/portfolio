package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

const indexFile = "index.html"

// spaHandler は SPA の静的ファイルを配信する。
// 実ファイルが存在すればそれを返し、存在しなければ index.html へ fallback する
// （クライアント側ルーティングの直リンクを 404 にしないため）。
//
// fs.FS を引数に取るのは、本番の embed.FS と単体テストの fstest.MapFS を
// 差し替え可能にするため。
func spaHandler(assets fs.FS) http.Handler {
	files := http.FileServerFS(assets)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !exists(assets, r.URL.Path) {
			serveIndex(w, r, assets)
			return
		}
		files.ServeHTTP(w, r)
	})
}

// exists は要求されたパスが実ファイルとして存在するかを返す。
func exists(assets fs.FS, urlPath string) bool {
	name := strings.TrimPrefix(path.Clean(urlPath), "/")
	if name == "" || name == "." {
		return false
	}
	info, err := fs.Stat(assets, name)
	return err == nil && !info.IsDir()
}

func serveIndex(w http.ResponseWriter, r *http.Request, assets fs.FS) {
	body, err := fs.ReadFile(assets, indexFile)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}
