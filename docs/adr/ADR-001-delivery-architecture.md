# ADR-001: 配信構成に Go 単一バイナリ + go:embed を採用する

- ステータス: Accepted
- 日付: 2026-08-19
- 関連: [ADR-002](./ADR-002-gitops-delivery.md), [前提](../assumptions-20260819.md)

## コンテキスト

Svelte 5 の SPA（ビルド成果物は完全に静的）と Go の API を、オンプレ VM 上の k3s / k3d にホスティングする。DB は持たず、HPA による水平 Pod スケールを機能させることが要件に含まれる。

リポジトリには既に `backend/static/` という空ディレクトリと、compose に「将来 go:embed で frontend/dist を参照できるように」というコメントが存在しており、当初からこの方向が想定されていた。

## 決定

Go バイナリに SPA を `//go:embed all:static` で埋め込み、**単一プロセスが静的配信と API の両方を担う**。本番イメージは multi-stage ビルドで distroless nonroot 上に置く。

## 理由

1. **版ずれが原理的に起きない** — フロントとバックが 1 つのアーティファクトに固定される
2. **HPA のスケール対象が 1 つ** — 2 コンテナ構成だと Deployment ごとに HPA を持つか、静的配信のために API ごとスケールする歪みを抱える
3. **完全ステートレス** — DB 無し・共有ストレージ無し・セッションアフィニティ不要で、HPA の前提を無条件に満たす
4. **同一オリジン** — 本番で CORS の考慮が不要になる
5. **オンプレ VM への負荷が小さい** — distroless で 15MB 程度に収まる

## 実装上の制約（実測で確認済み）

`//go:embed` は対象ディレクトリに埋め込み可能なファイルが 1 つも無いとコンパイルエラーになる。埋め込み先ディレクトリは git 未追跡の空ディレクトリになりがちで、クリーンチェックアウトでは存在すらしない。

検証結果:

| 条件 | 結果 |
|---|---|
| `//go:embed static` + 空ディレクトリ | `cannot embed directory static: contains no embeddable files` |
| `//go:embed static` + `.gitkeep` のみ | 同上（`.` 始まりは除外されるため） |
| `//go:embed all:static` + `.gitkeep` のみ | **成功** |
| `//go:embed static` + 通常ファイルあり | 成功 |

したがって **`all:` プレフィクスを付けたうえで `.gitkeep` をコミットする**。ダミーの `index.html` を置く必要はない。

### 配置

`//go:embed` のパスは埋め込む Go ファイルからの相対で、親ディレクトリを遡れない。当初案の `backend/static/` は `internal/server` から参照できないため、専用パッケージ `backend/internal/assets/` を作り、その配下に `static/` を置く。埋め込みの詳細を 1 パッケージに閉じ込められる利点もある。

## Vite の outDir を変更しない理由

`frontend/dist` のままにし、本番 Dockerfile 側で `static/` へ COPY する。outDir を `../backend/static` に向けると、開発中の `bun run build` が backend のツリーを汚し、`dist` を対象にした `.gitignore` とも噛み合わなくなる。

## 結果

- 静的ファイルのみの修正でも Go バイナリの再ビルドが必要
- 静的配信と API を個別にスケールできない
- gzip / brotli と `Cache-Control` は自前実装が必要（ETag と Range は標準実装で対応済み）

## 後日の変更

**2026-09-02**: デプロイ先が変わった。上記の決定（Go 単一バイナリ + `go:embed`）は維持する。

- コンテキストにある「オンプレ VM 上の k3s / k3d」は、**Talos Linux の単一ノード Kubernetes クラスタ**に変わった。クラスタは別リポジトリ `thin-k8s` が管理する
- 決定の根拠（完全ステートレス / HPA のスケール対象が 1 つ / 同一オリジン / distroless で小さい）はいずれもディストリビューションに依存しないため、決定は書き換えない
- HPA が依存する metrics-server は、k3s の同梱コンポーネントではなく `thin-k8s` の Talos machine config（`talos/patches/60-metrics-server.yaml`）が入れる
