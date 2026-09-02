# GitHub Actions 要件

デプロイ先は Talos Linux 上の単一ノード Kubernetes（`thin-k8s`）。
[ADR-002](adr/ADR-002-gitops-delivery.md) の pull 型 GitOps を実装する。

```
push to main
  → GitHub Actions（テスト → ビルド → GHCR push）
  → kustomize edit set image でタグ差し替え → コミット
  → Argo CD が検知して同期
```

## 前提（クラスタ側で用意済み）

| 項目 | 値 |
|---|---|
| Ingress Controller | Traefik（`ingressClassName: traefik`） |
| 外部公開 | Cloudflare Tunnel → Traefik。`cyokozai.net` |
| namespace | `portfolio` |
| imagePullSecret | `ghcr-pull`（`portfolio` namespace） |
| Argo CD | `argocd` namespace。`deploy/` を同期 |

## リポジトリ構成

```
backend/          Go 1.27。cmd/server, internal/
frontend/         Svelte + Vite。bun でビルド
container/
  Dockerfile      本番用。マルチステージで SPA を Go に埋め込む
deploy/           Kustomize。Argo CD の同期対象
  argocd/
    application.yaml
```

---

## W1: CI（PR 時）

**トリガー**: `pull_request` → `main`

**目的**: マージ前に壊れていないことを保証する。イメージは作らない。

### ジョブ

| # | 内容 | コマンド |
|---|---|---|
| 1 | Go のテスト | `cd backend && go test ./...` |
| 2 | Go の静的解析 | `go vet ./...` |
| 3 | フロントのビルド確認 | `cd frontend && bun install --frozen-lockfile && bun run build` |
| 4 | イメージのビルド確認 | `docker build -f container/Dockerfile .`（push しない） |

### 要件

- **Go 1.27** を使う（`backend/go.mod` に合わせる）
- **bun** を使う（`frontend/bun.lock` があるため npm/yarn は使わない）
- `--frozen-lockfile` でロックファイルとの不一致を検出する
- Go modules と bun の依存をキャッシュする
- 4 は 1〜3 が通ってから実行する（無駄なビルドを避ける）

---

## W2: CD（main への push 時）

**トリガー**: `push` → `main`

**目的**: イメージを GHCR へ publish し、`deploy/` のタグを更新する。

### ジョブ 1 — build & push

| 項目 | 要件 |
|---|---|
| イメージ名 | `ghcr.io/cyokozai/portfolio` |
| タグ | **コミット SHA**（`latest` は使わない） |
| 認証 | `GITHUB_TOKEN`（`packages: write` 権限） |
| Dockerfile | `container/Dockerfile` |
| ビルドコンテキスト | **リポジトリルート**（frontend と backend の両方が要る） |
| プラットフォーム | `linux/amd64` のみ |
| キャッシュ | GitHub Actions cache（`type=gha`） |

> ⚠ **`latest` タグを使わない理由**: `deploy/kustomization.yaml` のタグが変わらないと
> git に差分が出ず、Argo CD が同期しない。SHA タグが必須。

> ⚠ **プラットフォーム**: ノードは Intel i3（amd64）。arm64 は不要。
> ビルドするとその分だけ時間が延びる。

### ジョブ 2 — マニフェスト更新

ジョブ 1 の成功後に実行する。

1. `kustomize edit set image ghcr.io/cyokozai/portfolio=ghcr.io/cyokozai/portfolio:<SHA>` を `deploy/` で実行
2. 差分があれば commit して `main` へ push

**要件**:

- コミットメッセージに `[skip ci]` を含める（**再帰的な実行を防ぐ**）
- `contents: write` 権限が要る
- 差分が無い場合は commit しない
- コミット作者は `github-actions[bot]` を使う

> 🔴 **無限ループへの注意**: このジョブは `main` に push するため、
> `[skip ci]` が無いと W2 が自分自身を再トリガーし続ける。

---

## 権限

最小権限で設定する。

```yaml
permissions:
  contents: write     # マニフェスト更新の push
  packages: write     # GHCR への push
```

`GITHUB_TOKEN` で足りるため、**PAT の作成は不要**。

## Secrets

**追加の Secret は不要。** GHCR への push は `GITHUB_TOKEN` で行える。

クラスタ側の pull 用 PAT（`ghcr-pull`）は Actions とは無関係で、
クラスタに直接登録済み。

---

## 判断が必要な項目

実装前に決めること。

| # | 論点 | 選択肢 |
|---|---|---|
| 1 | GHCR を private にする場合の可視性設定 | GitHub の Packages 設定で明示的に private にする。リポジトリが public でもパッケージは独立して設定できる |
| 2 | フロントのテスト | 現状 test スクリプトが無い。追加するか、当面ビルド確認のみとするか |
| 3 | Lint | `golangci-lint` を入れるか、`go vet` のみとするか |
| 4 | ブランチ保護 | W1 の成功を必須にするか |

---

## 検証項目

実装後に確認すること。

- [ ] PR で W1 が走り、テストが失敗すると PR がブロックされる
- [ ] main への push で W2 が走り、GHCR に SHA タグのイメージが出る
- [ ] `deploy/kustomization.yaml` の `newTag` が SHA に更新される
- [ ] そのコミットが **W2 を再トリガーしない**（`[skip ci]` が効いている）
- [ ] Argo CD が差分を検知し、自動同期する
- [ ] Pod が新しいイメージで起動する（`kubectl -n portfolio get pod -o jsonpath='{..image}'`）
- [ ] `https://cyokozai.net` がサイトを返す
