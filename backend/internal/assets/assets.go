// Package assets はビルド済みの SPA をバイナリへ埋め込む。
//
// static/ の中身は本番イメージのビルド時に frontend/dist から COPY される。
// 開発時は空（.gitkeep のみ）で、SPA は Vite 開発サーバが配信する。
package assets

import (
	"embed"
	"io/fs"
)

// all: を付けているのは、static/ が .gitkeep しか持たない状態でも
// コンパイルを通すため。付けないと "." 始まりのファイルが除外され、
// "contains no embeddable files" でビルドできない。
//
//go:embed all:static
var embedded embed.FS

// FS は埋め込んだ SPA を、static/ を root として剥がした形で返す。
func FS() fs.FS {
	sub, err := fs.Sub(embedded, "static")
	if err != nil {
		panic(err) // embed 対象が固定なので、ここに来るのはビルド不整合のみ
	}
	return sub
}
