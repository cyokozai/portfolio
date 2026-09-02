# 実装計画

前提: [assumptions-20260819.md](./assumptions-20260819.md) / 設計: [architecture-decisions.md](./architecture-decisions.md)

## 方針

- **t_wada 流 TDD**（Red → Green → Refactor）を backend のロジックに適用する。マニフェストや Dockerfile はテスト対象外だが、各フェーズに明示的な受け入れ確認を置く
- **各フェーズの終わりで必ず green + コミット**。フェーズをまたいだ未完成状態を残さない
- **縦に 1 本通すことを優先**する。API を作り込む前に「ブラウザ → Vite proxy → Go → JSON」を最短で開通させ、その後に横へ広げる
- 依存の順に並べてある。Phase 0 を飛ばすと Phase 2 で必ず詰まる

## 進捗（2026-09-02 時点）

| # | フェーズ | 状態 |
|---|---|---|
| 0 | 開発基盤の統合 | **完了**（コンテナ間疎通を確認） |
| 1 | backend HTTP 骨格 | **完了**（TDD、SIGTERM の停止順序まで実プロセスで検証） |
| 2 | frontend 接続 | **完了**（ブラウザでの視覚確認のみ未） |
| 3 | 本番イメージ | **完了**（7.95MB / distroless / nonroot、`docker stop` で ExitCode 0） |
| 4 | k8s マニフェスト | **作成完了**（`kustomize build` 検証済み）。クラスタへの同期は未 |
| 5 | CI/CD | **実装済み**（PR 提出済み） |
| 6 | platform（helmfile） | **`thin-k8s` へ委譲**。portfolio 側の `platform/` は削除した |
| 7 | 公開（Application） | **実装済み**（`deploy/argocd/application.yaml`、PR 提出済み） |

### デプロイ先の変更（2026-09-02）

デプロイ先が VM 上の k3s から、**Talos Linux の単一ノードクラスタ**へ変わった。クラスタは別リポジトリ `thin-k8s` が管理する。責務の分界は次のとおり。

| 対象 | 管理者 |
|---|---|
| Traefik / Longhorn / cloudflared / Argo CD 本体 | `thin-k8s` の core 層（`helmfile/core/helmfile.yaml`） |
| metrics-server | `thin-k8s` の Talos machine config（`talos/patches/60-metrics-server.yaml`） |
| アプリ本体（`deploy/`） | portfolio |
| Argo CD の `Application` | portfolio（`deploy/argocd/application.yaml`） |

Ingress Controller は **Traefik**（namespace `traefik`、`isDefaultClass: true`）。公開ドメインは `cyokozai.net` で、Cloudflare Tunnel 経由。

### 次のアクション

1. **Phase 4/7 の実クラスタ確認**: Argo CD が `deploy/` を同期し、Pod が Running・Ingress 経由で `cyokozai.net` が 200 を返すこと
2. **HPA の確認**: `kubectl get hpa` の TARGETS が数値になること（metrics-server は Talos machine config が入れる）
3. **CI/CD の一巡確認**: `main` への push で GHCR へ push → タグ書き戻しコミット 1 つ → Argo CD が同期、で止まること
4. **`main` へのマージ**: GitHub Actions と Argo CD は `main` を見るため、作業ブランチのままでは反映されない

### 積み残し

- ブラウザでの視覚確認（Chrome 拡張が未接続のため未実施）
- HPA の実発火確認（クラスタ上で負荷をかけて行う）

---

## 全体像

| # | フェーズ | 主眼 | 目安 |
|---|---|---|---|
| 0 | 開発基盤の統合 | compose 1 本化。コンテナ間通信の開通 | 30 分 |
| 1 | backend HTTP 骨格 | readiness / SPA fallback / router / graceful shutdown（TDD） | 2〜3 時間 |
| 2 | frontend 接続 | Vite proxy、`/api/profile` 開通、死んだコードの掃除 | 1〜2 時間 |
| 3 | 本番イメージ | multi-stage → distroless 単一バイナリ | 1 時間 |
| 4 | k8s マニフェスト | Deployment / Service / HPA / Ingress。**検証は Talos クラスタ上で行う** | 1〜2 時間 |
| 5 | CI/CD | GitHub Actions（test / build / push / タグ差し替え） | 1〜2 時間 |
| 6 | platform（helmfile） | **`thin-k8s` へ委譲**（Argo CD 本体・cloudflared・Traefik は core 層が持つ） | ― |
| 7 | 公開 | Argo CD へ Application 登録、トンネル疎通確認 | 1 時間 |

**実行環境の方針（2026-08-22 決定 / 2026-09-02 更新）**: 実装・検証はすべてクラスタ側（Talos ノードと `thin-k8s` の作業環境）で行い、**Mac には何もインストールしない**。当初計画にあった Mac への軽量 k8s 導入は取りやめ。Mac 側で必要な確認は、使い捨てコンテナ（Docker は導入済み）で行う。

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

Talos クラスタに対して実施する。**手動 `kubectl apply` は運用手段に含めない**（[ADR-002](./adr/ADR-002-gitops-delivery.md)）ため、適用は Argo CD 経由で行い、確認はその結果に対して行う。

```bash
kustomize build deploy/ | kubectl diff -f -   # 同期前の差分確認（読み取りのみ）
argocd app sync portfolio                     # または Argo CD の自動同期を待つ
kubectl -n portfolio get hpa                  # TARGETS が <unknown> でなく数値になること
kubectl -n portfolio rollout restart deploy/portfolio && kubectl -n portfolio get pods -w   # 502 が出ないこと
```

metrics-server は `thin-k8s` の Talos machine config（`talos/patches/60-metrics-server.yaml`）が入れる。Talos には k3s のような同梱コンポーネントが無いため、**追加導入が不要という前提は成り立たない**。helmfile には置かず machine config 側に寄せてある。HPA の発火は負荷をかけて確認する。

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

## Phase 6: platform（helmfile）→ `thin-k8s` へ委譲（2026-09-02）

**このフェーズは portfolio の担当から外れた。** クラスタの土台は Talos クラスタ側リポジトリ `thin-k8s` が持つ。portfolio 側にあった `platform/`（helmfile と values）は削除した。

同じものを入れる helmfile が 2 つ並ぶと、`helmfile apply` で 2 つめの Argo CD（portfolio 側 10.6.0 / `thin-k8s` 側 10.4.2）と 2 つめの cloudflared を建てにいくため、重複を残さず削除する判断にした。

### 移った先

| 旧（portfolio） | 新（`thin-k8s`） |
|---|---|
| `platform/helmfile.yaml` の `argo/argo-cd` | `helmfile/core/helmfile.yaml`（Argo CD 10.4.2） |
| `platform/helmfile.yaml` の `cloudflare/cloudflare-tunnel` | `helmfile/core/helmfile.yaml`（chart 0.3.2） |
| `platform/values/argocd.yaml.gotmpl` | `helmfile/core/values/` 配下 |
| `platform/values/cloudflare-tunnel.yaml.gotmpl` | `helmfile/core/values/` 配下 |
| `platform/README.md`（トンネル作成手順） | `thin-k8s` の `SETUP.md` / `OPERATIONS.md` |
| （新規） | Traefik 41.4.0 / Longhorn、metrics-server は `talos/patches/60-metrics-server.yaml` |

### 責務の分界（変わっていない部分）

| レイヤ | 管理者 | 対象 |
|---|---|---|
| day-0 | `thin-k8s`（helmfile / machine config） | Traefik、Longhorn、cloudflared、Argo CD 本体、metrics-server |
| day-2 | Argo CD | アプリ本体（`deploy/`） |

Argo CD が Argo CD 自身を管理する循環を避けるため、両者の対象は重ねない。この原則は担当が移っても変わらない。

### portfolio 側に残る前提

- トンネル認証情報の Secret は帯域外で登録する（リポジトリには置かない）
- Argo CD の UI は既定で公開しない。GitOps の制御面をインターネットに晒さないため、`port-forward` で見る
- Ingress は `ingressClassName: traefik` を前提にする（`deploy/ingress.yaml`）

---

## Phase 7: 公開

### 作業

- `deploy/argocd/application.yaml` — `repoURL` は本リポジトリ、`path: deploy`、`targetRevision: main`、`syncPolicy.automated` に `prune: true` / `selfHeal: true`
- Argo CD への初回登録は**手動 bootstrap**（GitOps の対象外）。Argo CD 本体は `thin-k8s` が入れる

### 受け入れ確認

main へ push → cd がタグを書き戻す → 数分以内に Pod のイメージが入れ替わる。`cyokozai.net` で HTTPS 公開され、`/about` の直リンクが 200 で返る

---

## このあと必要になったら作るもの

`arch-requirements` の残りの成果物は、現時点では過剰と判断して省いている。必要になった時点で作成する。

- `docs/slo.md` — 可用性 SLO を明示したくなったとき
- `docs/finops.md` — オンプレ VM の電力・回線コストを可視化したくなったとき
- `docs/runbook.md` — 障害対応手順を残したくなったとき
- `docs/prd.md` — 掲載コンテンツの要件を詰めるとき
