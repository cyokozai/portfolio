# 実装計画

前提: [assumptions-20260819.md](./assumptions-20260819.md) / 設計: [architecture-decisions.md](./architecture-decisions.md)

## 方針

- **t_wada 流 TDD**（Red → Green → Refactor）を backend のロジックに適用する。マニフェストや Dockerfile はテスト対象外だが、各フェーズに明示的な受け入れ確認を置く
- **各フェーズの終わりで必ず green + コミット**。フェーズをまたいだ未完成状態を残さない
- **縦に 1 本通すことを優先**する。API を作り込む前に「ブラウザ → Vite proxy → Go → JSON」を最短で開通させ、その後に横へ広げる
- 依存の順に並べてある。Phase 0 を飛ばすと Phase 2 で必ず詰まる

## 進捗（2026-08-26 時点）

| # | フェーズ | 状態 |
|---|---|---|
| 0 | 開発基盤の統合 | **完了**（コンテナ間疎通を確認） |
| 1 | backend HTTP 骨格 | **完了**（TDD、SIGTERM の停止順序まで実プロセスで検証） |
| 2 | frontend 接続 | **完了**（ブラウザでの視覚確認のみ未） |
| 3 | 本番イメージ | **完了**（7.95MB / distroless / nonroot、`docker stop` で ExitCode 0） |
| 4 | k8s マニフェスト | **作成完了**（`kustomize build` 検証済み）。クラスタへの適用は未 |
| 5 | CI/CD | **未着手** ← 次はここ |
| 6 | platform（helmfile） | **作成完了**（`helmfile template` EXIT=0）。クラスタへの適用は未 |
| 7 | 公開 | **未着手**（`argocd/application.yaml` が未作成） |

### 次のアクション

1. **Phase 5**: `.github/workflows/ci.yaml` と `cd.yaml` を作成（Mac 側で作業可能）
2. **Phase 7**: `argocd/application.yaml` を作成（Mac 側で作業可能）
3. **`CHANGE_ME` の置換**: `platform/helmfile.yaml`（domain / tunnelId / accountId）と `deploy/ingress.yaml`（host）。domain は 2 ファイルに出てくるので揃える
4. **VM 上での適用**（ユーザー作業）: トンネル作成 → Secret 登録 → `helmfile apply` → `kubectl apply -k deploy/` → ArgoCD へ Application 登録
5. **push**: ローカルが `origin/feat/backend` より 2 コミット先行。GitHub Actions と ArgoCD は `main` を見るため、`main` へのマージも必要

### 積み残し

- ブラウザでの視覚確認（Chrome 拡張が未接続のため未実施）
- HPA の実発火確認（VM 上で負荷をかけて行う）

---

## 全体像

| # | フェーズ | 主眼 | 目安 |
|---|---|---|---|
| 0 | 開発基盤の統合 | compose 1 本化。コンテナ間通信の開通 | 30 分 |
| 1 | backend HTTP 骨格 | readiness / SPA fallback / router / graceful shutdown（TDD） | 2〜3 時間 |
| 2 | frontend 接続 | Vite proxy、`/api/profile` 開通、死んだコードの掃除 | 1〜2 時間 |
| 3 | 本番イメージ | multi-stage → distroless 単一バイナリ | 1 時間 |
| 4 | k8s マニフェスト | Deployment / Service / HPA。**検証は VM 上の k3s で行う** | 1〜2 時間 |
| 5 | CI/CD | GitHub Actions（test / build / push / タグ差し替え） | 1〜2 時間 |
| 6 | platform（helmfile） | ArgoCD 本体と cloudflared を VM の k3s へ導入 | 1 時間 |
| 7 | 公開 | ArgoCD へ Application 登録、トンネル疎通確認 | 1 時間 |

**実行環境の方針（2026-08-22 決定）**: 実装・検証はすべて VM 上で行い、**Mac には何もインストールしない**。当初計画にあった Mac への `k3d` 導入は取りやめ。Mac 側で必要な確認は、使い捨てコンテナ（Docker は導入済み）で行う。

---

## Phase 0: 開発基盤の統合

現状 2 つの compose が別 project 名（`portfolio-frontend` / `portfolio-backend`）で動いており、**ネットワークが分離している**。Phase 2 でフロントから API を叩いた瞬間に通信できないため、先に潰す。

### 作業

1. ルートに `compose.yaml` を新規作成し、`frontend` / `backend` の 2 サービスを同一ネットワークに置く
   - build context はそれぞれ既存の `frontend/container/Dockerfile`・`backend/container/Dockerfile` を流用する
   - ports: `5173:5173`（Vite）、`8080:8080`（Go）
   - volumes: リポジトリルートを `/workspace` へ。`gomodcache` と `node_modules` は named volume を維持
2. `.devcontainer/backend/devcontainer.json`・`.devcontainer/frontend/devcontainer.json` の `dockerComposeFile` を `../../compose.yaml` に差し替える
3. 旧 `frontend/container/compose.yaml`・`backend/container/compose.yaml` を削除（Dockerfile は残す）

### 注意

- **既存の named volume は捨てることになる**。project 名が変わると別 volume が作られるため、`gomodcache` と `node_modules` は再取得になる
- `updateRemoteUserUID: false` は維持する（macOS で uid が書き換わると gomodcache の所有者と食い違う）

### 受け入れ確認

```bash
docker compose exec frontend curl -sf http://backend:8080/health
docker compose exec backend  curl -sfI http://frontend:5173
```

---

## Phase 1: backend HTTP 骨格（TDD）

`/health` は実装済み。ここに readiness・SPA 配信・ルーティング・graceful shutdown を足す。

### 1-1. readiness（`internal/server/ready.go`）

- **Red**: `TestReadyHandler` — 初期状態で 200、`markShuttingDown()` 呼び出し後に 503
- **Green**: `atomic.Bool` を持つ最小実装

### 1-2. SPA fallback（`internal/server/static.go`）

**設計判断**: ハンドラのシグネチャを `spaHandler(fsys fs.FS) http.Handler` にする。`embed.FS` を直接参照すると単体テストで実ファイルが必要になるため、**テストでは `fstest.MapFS` を注入する**。`embed.FS` の生成は本番コードの 1 箇所だけに閉じ込める。

- **Red**: `TestSPAHandler`
  - 実ファイルが存在するパス（`/assets/app.js`）→ その中身を返す
  - 存在しないパス（`/about`）→ 200 + `index.html` の中身
- **Green**: `http.FS` + fallback 分岐

### 1-3. router（`internal/server/router.go`）

- **Red**: `TestRouter` — `/health` 200 / `/ready` 200 / `/api/unknown` **404**（`index.html` に fallback しないこと）/ `/` と `/about` は `index.html`
- **Green**: `http.ServeMux`（Go 1.22 以降のパターンマッチが使える）で `/api/` を明示登録し、fallback から除外する

### 1-4. graceful shutdown（`cmd/server/main.go`）

**HPA のスケールインとローリング更新で Pod が日常的に落ちるため必須。** SIGTERM → readiness を 503 に反転 → 数秒待つ → `http.Server.Shutdown()` の順序を守る。この順序を守らないと Service がまだ Pod を宛先に含んでいる間に接続が切れて 502 になる。

- **Red**: 起動処理を `run(ctx context.Context) error` に切り出し、ctx キャンセルで readiness が false になり `Shutdown` が完了することを検証
- **Green**: `signal.NotifyContext` + `Shutdown`
- `main()` は `run()` を呼ぶだけの薄い関数にする

### 1-5. embed の下地

**配置**: `//go:embed` は親ディレクトリを遡れないため、埋め込み先は埋め込む Go ファイルと同じ階層に置く必要がある。`backend/static/` は `internal/server` から参照できないので、専用パッケージ `backend/internal/assets/` を作り、その配下に `static/` を置く。

- `backend/internal/assets/assets.go` — `//go:embed all:static` と `FS()`（`fs.Sub` で `static/` を root として剥がす）
- `backend/internal/assets/static/.gitkeep` — **コミットする**

`all:` が無いと `.` 始まりのファイルは除外され `contains no embeddable files` でコンパイルできない（実測確認済み）。

### 受け入れ確認

`go test ./...` green / `go vet ./...` 無警告 / `gofmt -l .` 空

---

## Phase 2: frontend 接続

### 作業

1. `frontend/vite.config.js`
   - `server.host: '0.0.0.0'` — 既定では localhost にのみ bind するため、Docker のポート公開越しにホストから到達できない
   - `server.proxy: { '/api': 'http://backend:8080' }`
2. `/api/profile` を TDD で 1 本追加（氏名・肩書・所属を JSON で返す。データは Go の構造体リテラルで十分）
3. `App.svelte` から `fetch('/api/profile')` してヒーロー部を描画する
4. **死んだコードの掃除**（`git grep` で確認済み・いずれもマークアップ側から参照されていない）
   - 未使用の宣言: `onVisible()`、`mottosVisible`
   - 未参照 CSS: `.rule--yellow` `.hero-mottos` `.motto` `#mottos` `.cyan-block` `.block-motto` `.img-label`
   - 空要素: `<p class="hero-tagline"></p>`、`hero-affil` の空 `<li>`

### 受け入れ確認

ホストのブラウザで `http://localhost:5173` を開き、ヒーロー部が API 由来のデータで描画される。DevTools の Network で `/api/profile` が 200。

---

## Phase 3: 本番イメージ

### 作業

**配置**: リポジトリルートに `container/Dockerfile` を新規作成する。frontend と backend の両方のソースが要るためビルドコンテキストがリポジトリルートになり、`backend/container/` 配下には置けない。ルートに `.dockerignore` も作る。

3 ステージ構成:

| stage | ベース | 処理 |
|---|---|---|
| web | `oven/bun:1` | `bun install --frozen-lockfile` → `bun run build` → `frontend/dist` |
| build | `golang:1.26` | web の成果物を `backend/internal/assets/static/` へ COPY → `go build -trimpath -ldflags="-s -w"` |
| final | `gcr.io/distroless/static-debian12:nonroot` | バイナリのみ COPY、`USER nonroot`、`EXPOSE 8080` |

### 受け入れ確認

```bash
docker build -f container/Dockerfile -t portfolio:dev .
docker run --rm -p 8080:8080 portfolio:dev
curl -sf localhost:8080/api/profile      # JSON
curl -sf localhost:8080/ | head          # SPA
curl -so /dev/null -w '%{http_code}\n' localhost:8080/about   # 200（fallback）
curl -so /dev/null -w '%{http_code}\n' localhost:8080/api/x   # 404
docker images portfolio:dev              # 20MB 未満を期待
```

**実績（2026-08-22）**: イメージ 7.95MB / `User=nonroot:nonroot` / シェル・`/bin/ls` とも不在 / `docker stop` で `/ready` 503・`/health` 200 を保ったまま ExitCode 0 で正常終了。

---

## Phase 4: k8s マニフェスト

### 作業

`deploy/` に単一セットで置く（overlay は作らない）。

- `deployment.yaml`
  - `replicas: 1`
  - `resources.requests`: `cpu: 50m` / `memory: 64Mi` — **HPA が使用率を計算するために requests は必須**。過大に置くと HPA が永久に発火しない
  - `resources.limits`: `memory: 128Mi`。**CPU limit は付けない**（throttling で HPA の判断が歪むため）
  - probes: liveness → `/health`、readiness → `/ready`
  - `securityContext`: `runAsNonRoot: true`、`readOnlyRootFilesystem: true`（静的埋め込みなので書き込み不要）、`allowPrivilegeEscalation: false`
- `service.yaml` — ClusterIP、8080
- `hpa.yaml` — `autoscaling/v2`、`minReplicas: 1`、`maxReplicas: 3`、CPU `averageUtilization: 70`
- `kustomization.yaml` — 上記を列挙。`images:` に GHCR のイメージを記述（Phase 5 の差し替え対象）

### 受け入れ確認

VM 上の k3s に対して実施する。

```bash
kubectl apply -k deploy/
kubectl get hpa                # TARGETS が <unknown> でなく数値になること
kubectl rollout restart deploy/portfolio && kubectl get pods -w   # 502 が出ないこと
```

metrics-server は k3s に同梱されているため追加導入は不要（`--disable` で外せる packaged component の一覧に含まれることを確認済み）。HPA の発火は負荷をかけて確認する。

---

## Phase 5: CI/CD（GitHub Actions）

### 作業

- `.github/workflows/ci.yaml` — `pull_request` と `push: main`
  - backend job: `go test ./...` / `go vet ./...` / `gofmt -l` が空であること
  - frontend job: `bun install --frozen-lockfile` / `bun run build`
- `.github/workflows/cd.yaml` — `push: main`
  - GHCR へ `ghcr.io/cyokozai/portfolio:<git sha>` を push（`permissions: packages: write`）
  - `kustomize edit set image` でタグを差し替え、リポジトリへコミットを戻す（`permissions: contents: write`）

### 注意

- **両ワークフローに `paths-ignore: ['deploy/**']` を必ず設定する。** 忘れると bot のコミットが CI を再起動させ、無限に自走する
- 書き戻しコミットのメッセージに `[skip ci]` を付けて二重の保険にする
- イメージタグは `latest` ではなく git sha にする。ArgoCD が差分を検知できず同期しなくなるため

### 受け入れ確認

PR で ci が緑 / main への push で cd がタグ更新コミットを 1 つだけ作り、そこで停止する

---

## Phase 6: platform（helmfile）

クラスタの土台を helmfile でまとめて入れる。**作成済み**（`platform/`）。

### 責務の分界

| レイヤ | 管理者 | 対象 |
|---|---|---|
| day-0 | helmfile | ArgoCD 本体、cloudflared |
| day-2 | ArgoCD | アプリ本体（`deploy/`） |

ArgoCD が ArgoCD 自身を管理する循環を避けるため、両者の対象は重ねない。

### 構成

- `platform/helmfile.yaml` — `argo/argo-cd` 10.4.0、`cloudflare/cloudflare-tunnel` 0.3.2
- `platform/values/argocd.yaml.gotmpl` — 単一ノード向けに縮小。`applicationSet` / `notifications` / `dex` は無効
- `platform/values/cloudflare-tunnel.yaml.gotmpl` — ingress ルールを Git に残す（remotely-managed 方式は採らない）
- `platform/README.md` — トンネル作成から適用までの手順

### 注意

- トンネル認証情報の Secret は**帯域外で登録する**。chart の `secretName` を指定して chart 側に Secret を作らせない。キー名は `credentials.json` 固定（chart が `/etc/cloudflared/creds/credentials.json` を読むため）
- **ArgoCD の UI は既定で公開しない。** GitOps の制御面をインターネットに晒さないため、`port-forward` で見る

### 受け入れ確認

`helmfile -f platform/helmfile.yaml apply` 後に ArgoCD と cloudflared の Pod が Running

---

## Phase 7: 公開

### 作業

- `argocd/application.yaml` — `repoURL` は本リポジトリ、`path: deploy`、`targetRevision: main`、`syncPolicy.automated` に `prune: true` / `selfHeal: true`
- ArgoCD への初回登録は**手動 bootstrap**（GitOps の対象外）

### 受け入れ確認

main へ push → cd がタグを書き戻す → 数分以内に Pod のイメージが入れ替わる。独自ドメインで HTTPS 公開され、`/about` の直リンクが 200 で返る

---

## このあと必要になったら作るもの

`arch-requirements` の残りの成果物は、現時点では過剰と判断して省いている。必要になった時点で作成する。

- `docs/slo.md` — 可用性 SLO を明示したくなったとき
- `docs/finops.md` — オンプレ VM の電力・回線コストを可視化したくなったとき
- `docs/runbook.md` — 障害対応手順を残したくなったとき
- `docs/prd.md` — 掲載コンテンツの要件を詰めるとき
