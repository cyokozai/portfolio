# platform — クラスタの土台（day-0）

VM 上の k3s に、アプリより下のレイヤを入れるための helmfile。

## 責務の分界

| レイヤ | 管理者 | 対象 |
|---|---|---|
| day-0 | **helmfile**（このディレクトリ） | ArgoCD 本体、cloudflared |
| day-2 | **ArgoCD** | アプリ本体（`deploy/`） |

ArgoCD が ArgoCD 自身を管理する循環を避けるため、両者の対象は重ねない。

`metrics-server` は k3s に同梱されているため入れていない（`--disable` で外せる packaged component の一覧に含まれる）。HPA はこれに依存するので、無効化している場合は `helmfile.yaml` 末尾のコメントアウトされたリリースを有効化する。

## 前提

VM 上に `helm` と `helmfile` が入っていること。Mac 側には何も入れない。

## 手順

### 1. トンネルを作成する（Cloudflare 側の作業）

```bash
cloudflared tunnel login
cloudflared tunnel create portfolio
```

`~/.cloudflared/<TUNNEL_ID>.json` が生成される。**このファイルはリポジトリに置かない。**

### 2. 認証情報を Secret としてクラスタへ登録する

chart は `/etc/cloudflared/creds/credentials.json` を読むため、キー名を `credentials.json` にする必要がある。

```bash
kubectl create namespace cloudflare

kubectl create secret generic cloudflare-tunnel-credentials \
  --namespace cloudflare \
  --from-file=credentials.json=$HOME/.cloudflared/<TUNNEL_ID>.json
```

この Secret は GitOps の対象外（帯域外で登録する）。`helmfile.yaml` の `tunnelSecretName` がこの名前を参照している。

### 3. DNS を向ける

```bash
cloudflared tunnel route dns portfolio <DOMAIN>
```

### 4. 値を埋める

`helmfile.yaml` の `environments.default.values` にある `CHANGE_ME` を置き換える。

| キー | 取得元 |
|---|---|
| `domain` | 公開するホスト名 |
| `tunnelId` | `cloudflared tunnel list` |
| `accountId` | Cloudflare ダッシュボード |

### 5. 適用する

```bash
helmfile -f platform/helmfile.yaml diff    # 差分を確認
helmfile -f platform/helmfile.yaml apply
```

### 6. ArgoCD に Application を登録する（初回のみ）

```bash
kubectl apply -f argocd/application.yaml
```

以降、`main` への push は ArgoCD が自動で同期する。

## 公開方針

**外に出すのは Web サイトのみ。** ArgoCD をはじめ、それ以外のリソースはトンネルにも Ingress にも載せない。

ArgoCD の UI は port-forward で参照する。

```bash
kubectl -n argocd port-forward svc/argocd-server 8080:80
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath='{.data.password}' | base64 -d
```

## チャートのバージョン

| チャート | バージョン | アプリ |
|---|---|---|
| `argo/argo-cd` | 10.4.0 | Argo CD v3.5.1 |
| `cloudflare/cloudflare-tunnel` | 0.3.2 | cloudflared 2024.8.3 |

`cloudflare-tunnel` の appVersion は古いので、必要なら `image.tag` で cloudflared のバージョンを上書きする。
